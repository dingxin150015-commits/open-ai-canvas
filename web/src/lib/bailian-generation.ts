import { wan27VideoKind, happyHorse11VideoKind, type VideoCapabilityConfig } from "./model-capabilities";
import type { AiConfig } from "@/stores/use-config-store";

export function declaredGenerationContract(config: Pick<AiConfig, "channels" | "model">, model = config.model): string | undefined {
    const separator = model.indexOf("::");
    const channelId = separator >= 0 ? model.slice(0, separator) : "";
    const name = separator >= 0 ? model.slice(separator + 2) : model;
    const channel = config.channels.find((item) => item.id === channelId) ?? config.channels.find((item) => item.models.includes(name));
    const cost = channel?.modelCosts?.find((item) => item.model === name);
    const direct = cost?.capabilityConfig?.generationContract;
    if (direct && bailianModelKind(direct)) return direct;
    const values = cost?.logicalCapabilitySpec?.options?.generationContract?.values;
    return values?.length === 1 && typeof values[0] === "string" && bailianModelKind(values[0]) ? values[0] : undefined;
}

export function bailianContractModel(config: Pick<AiConfig, "channels" | "model">, model = config.model) {
    return declaredGenerationContract(config, model) ?? model;
}

export function nativeImageBatchEnabled(config: AiConfig) {
    return (config.taskWorkflowProvider || "model") === "model" && config.generationOptions?.nativeBatch === true && bailianModelKind(bailianContractModel(config)) === "qwen" && (!config.generationOptions.contractKind || config.generationOptions.contractKind === "qwen");
}

export type BailianGenerationOptions = {
    contractKind?: string;
    seed?: string | number;
    negative_prompt?: string;
    prompt_extend?: boolean;
    prompt_extend_mode?: "direct" | "agent";
    enable_thinking?: boolean;
    watermark?: boolean;
    request_mode?: "sync" | "async" | "openai";
    nativeBatch?: boolean;
};

export function bailianReferencePrompt(prompt: string, model: string, images: { id: string }[], firstFrameID?: string) {
    const kind = bailianModelKind(model);
    if (kind !== "happyhorse-1.1-r2v" && kind !== "wan2.7-r2v") return prompt;
    let referenceIndex = 0;
    const labels = images.map((image) => (kind === "wan2.7-r2v" && image.id === firstFrameID ? "首帧" : kind.startsWith("happyhorse") ? `[Image ${++referenceIndex}]` : `Image ${++referenceIndex}`));
    const label = (type: string, index: string, fallback: string) => (type === "图片" ? (labels[Number(index) - 1] ?? fallback) : kind === "wan2.7-r2v" ? "Video " + index : fallback);
    const mapped = prompt.replace(/^(- )(@(图片|视频)(\d+))：参考(?:图片|视频) \d+$/gm, (_match, prefix: string, token: string, type: string, index: string) => prefix + token + "：" + label(type, index, token));
    return mapped.replace(/@(图片|视频)(\d+)(?!\d)/g, (token, type: string, index: string) => label(type, index, token));
}

export function validateBailianOptions(options: BailianGenerationOptions, model: string, imageCount: number) {
    if (!bailianModelKind(model)) return;
    const seed = options.seed === undefined ? "" : String(options.seed).trim();
    if (seed && (!/^\d+$/.test(seed) || Number(seed) > 2147483647)) throw new Error("随机种子须为 0–2147483647 的整数");
    if ([...(options.negative_prompt ?? "")].length > 500) throw new Error("负向提示词不能超过 500 字符");
    if (bailianModelKind(model) === "qwen" && imageCount > 0 && options.prompt_extend_mode === "agent") throw new Error("图像编辑不支持 Agent 扩写，请选择 Direct");
}

export function validateBailianFrames(model: string, images: { id: string }[], videos: number, audios: number, metadata: Record<string, unknown> = {}) {
    const kind = bailianModelKind(model);
    if (!kind || kind === "qwen") return;
    const first = metadata.videoStartFrameNodeId;
    const last = metadata.videoEndFrameNodeId;
    const roles = videoFrameRoles(model);
    if ((first && !roles.first) || (last && !roles.last)) throw new Error("当前模型不支持所选首尾帧角色，请清除旧设置");
    if ((first && !images.some((image) => image.id === first)) || (last && !images.some((image) => image.id === last))) throw new Error("指定首尾帧不在参考素材中，请重新选择");
    if (first && first === last) throw new Error("首帧和尾帧不能使用同一个素材槽位");
    if (kind === "wan2.7-i2v") {
        if (videos && (audios || images.length > 1 || first)) throw new Error("视频续写仅支持输入视频及可选尾帧，不支持驱动音频或另加首帧");
        if (!videos && last && !first) throw new Error("选择尾帧时请同时指定首帧");
    }
}

export function bailianModelKind(model: string) {
    const name = model
        .split("::")
        .at(-1)!
        .trim()
        .toLowerCase()
        .replace(/^models\//, "");
    if (name === "qwen-image-3.0-pro") return "qwen";
    if (happyHorse11VideoKind(name)) return `happyhorse-1.1-${happyHorse11VideoKind(name)}`;
    if (wan27VideoKind(name)) return `wan2.7-${wan27VideoKind(name)}`;
    return "";
}

export function optionalVideoPrompt(model: string, images: number, videos: number) {
    const kind = bailianModelKind(model);
    return (kind === "happyhorse-1.1-i2v" && images === 1 && videos === 0) || (kind === "wan2.7-i2v" && images + videos > 0);
}

export function bailianVideoOperation(model: string, videos: number, audios: number): string | undefined {
    const kind = bailianModelKind(model);
    if (kind.endsWith("-r2v")) return "reference_to_video";
    if (kind.endsWith("-i2v")) return videos > 0 && kind.startsWith("wan") ? "extend" : "image_to_video";
    if (kind.endsWith("-t2v")) return audios > 0 ? "audio_to_video" : "text_to_video";
    return undefined;
}

export function videoFrameRoles(model: string): { first: boolean; last: boolean } {
    const kind = bailianModelKind(model);
    if (!kind) return { first: true, last: true };
    return { first: kind === "happyhorse-1.1-i2v" || kind === "wan2.7-i2v" || kind === "wan2.7-r2v", last: kind === "wan2.7-i2v" };
}

export function videoAttachmentLimit(profile: VideoCapabilityConfig) {
    const refs = profile.references;
    const visualLimit = refs.maxVisualReferences && refs.maxVisualReferences > 0 ? refs.maxVisualReferences : refs.maxImages + refs.maxVideos;
    return Math.min(visualLimit, refs.maxImages + refs.maxVideos) + refs.maxAudios;
}

export function videoAttachmentError(profile: VideoCapabilityConfig, kinds: string[]) {
    const count = (kind: string) => kinds.filter((item) => item === kind).length;
    const refs = profile.references;
    if (count("image") > refs.maxImages) return `当前模型最多支持 ${refs.maxImages} 张图片`;
    if (count("video") > refs.maxVideos) return `当前模型最多支持 ${refs.maxVideos} 个视频`;
    if (count("audio") > refs.maxAudios) return `当前模型最多支持 ${refs.maxAudios} 个音频`;
    if (refs.maxVisualReferences && count("image") + count("video") > refs.maxVisualReferences) return `图片和视频合计最多 ${refs.maxVisualReferences} 个`;
    return "";
}
