import type { ModelProtocol } from "@/lib/model-protocols";

export type ModelCapabilityConfig = {
    version: number;
    text?: TextCapabilityConfig;
    image?: ImageCapabilityConfig;
    video?: VideoCapabilityConfig;
};

export type TextCapabilityConfig = {
    references: {
        promptMaxChars: number;
        maxImages: number;
        maxImageBytes: number;
        maxVideos: number;
        maxVideoBytes: number;
    };
};

export type ImageSizeParameter = "none" | "size" | "aspect_ratio";

export type ImageCapabilityConfig = {
    references: {
        promptMaxChars: number;
        maxImages: number;
        maxImageBytes: number;
        maskSupported: boolean;
    };
    size: {
        parameter: ImageSizeParameter;
        values: string[];
        default: string;
        allowCustom: boolean;
    };
    quality: {
        supported: boolean;
        values: string[];
        default: string;
    };
    transparentBackground: { supported: boolean; default: boolean };
    responseFormat: { supported: boolean };
    outputFormat: { supported: boolean };
    maxOutputs: number;
};

export type VideoCapabilityConfig = {
    references: {
        promptMaxChars: number;
        minVisualReferences?: number;
        maxVisualReferences?: number;
        minImages: number;
        maxImages: number;
        maxImageBytes: number;
        maxVideos: number;
        maxVideoBytes: number;
        maxVideoDurationSeconds: number;
        maxOutputDurationWithVideoSeconds?: number;
        maxAudios: number;
        maxAudioBytes: number;
        maxAudioDurationSeconds: number;
        audioMustFitOutput?: boolean;
    };
    duration: {
        selection: "range" | "enum";
        min?: number;
        max?: number;
        step?: number;
        values?: number[];
        default: number;
        smartSupported?: boolean;
    };
    ratios: string[];
    defaultRatio: string;
    resolutions: string[];
    defaultResolution: string;
    generateAudio: { supported: boolean; default: boolean };
    watermark: { supported: boolean; default: boolean };
    operations: string[];
    defaultOperation: string;
};

// 旧版本的“允许自定义”可能只保存了 `*`，前台需要用这组标准值恢复可选项。
export const STANDARD_IMAGE_SIZE_VALUES = [
    "1:1",
    "3:2",
    "2:3",
    "4:3",
    "3:4",
    "16:9",
    "21:9",
    "9:16",
    "1024x1024",
    "1536x1024",
    "1024x1536",
] as const;

export function normalizeCapabilityString(value: string) {
    const normalized = value.trim();
    return normalized.startsWith("string:") ? normalized.slice("string:".length) : normalized;
}

function normalizeCapabilityStrings(values: string[]) {
    return Array.from(new Set(values.map(normalizeCapabilityString)));
}

export function normalizeModelCapabilityConfig(config: ModelCapabilityConfig): ModelCapabilityConfig {
    return {
        ...config,
        image: config.image
            ? {
                  ...config.image,
                  size: {
                      ...config.image.size,
                      values: normalizeCapabilityStrings(config.image.size.values),
                      default: normalizeCapabilityString(config.image.size.default),
                  },
                  quality: {
                      ...config.image.quality,
                      values: normalizeCapabilityStrings(config.image.quality.values),
                      default: normalizeCapabilityString(config.image.quality.default),
                  },
              }
            : undefined,
        video: config.video
            ? {
                  ...config.video,
                  ratios: normalizeCapabilityStrings(config.video.ratios),
                  defaultRatio: normalizeCapabilityString(config.video.defaultRatio),
                  resolutions: normalizeCapabilityStrings(config.video.resolutions),
                  defaultResolution: normalizeCapabilityString(config.video.defaultResolution),
                  operations: normalizeCapabilityStrings(config.video.operations),
                  defaultOperation: normalizeCapabilityString(config.video.defaultOperation),
              }
            : undefined,
    };
}

// Keep explicit pixel presets for each resolution tier so the settings panel can
// switch between 1K, 2K and 4K without silently converting the requested ratio.
const defaultImageSizes = [
    "auto",
    "1:1",
    "3:2",
    "2:3",
    "4:3",
    "3:4",
    "16:9",
    "21:9",
    "9:16",
    "1024x1024",
    "1360x1024",
    "1024x1360",
    "1536x1024",
    "1024x1536",
    "1024x1280",
    "1280x1024",
    "2048x878",
    "1824x1024",
    "1024x1824",
    "2048x2048",
    "2304x1728",
    "1728x2304",
    "2496x1664",
    "1664x2496",
    "1792x2240",
    "2240x1792",
    "3136x1344",
    "2752x1536",
    "1536x2752",
    "2880x2880",
    "3264x2448",
    "2448x3264",
    "3504x2336",
    "2336x3504",
    "2560x3200",
    "3200x2560",
    "3808x1632",
    "3840x2160",
    "2160x3840",
];

export function defaultImageCapabilityConfig(protocol?: ModelProtocol, model = ""): ImageCapabilityConfig {
    const image: ImageCapabilityConfig = {
        references: { promptMaxChars: 32000, maxImages: 16, maxImageBytes: 30 * 1024 * 1024, maskSupported: true },
        size: { parameter: "size", values: [...defaultImageSizes], default: "1:1", allowCustom: true },
        quality: { supported: true, values: ["auto", "low", "medium", "high"], default: "auto" },
        transparentBackground: { supported: true, default: false },
        responseFormat: { supported: true },
        outputFormat: { supported: true },
        maxOutputs: 15,
    };
    if (protocol === "grok-image") {
        image.references.maxImages = 1;
        image.references.maskSupported = false;
        // grok2api / xAI Imagine：size→aspect_ratio，quality→resolution(1k/2k)。
        image.size = {
            parameter: "aspect_ratio",
            values: ["1:1", "3:4", "4:3", "9:16", "16:9", "2:3", "3:2"],
            default: "1:1",
            allowCustom: false,
        };
        image.quality = { supported: true, values: ["1k", "2k"], default: "2k" };
        image.transparentBackground = { supported: false, default: false };
        image.responseFormat = { supported: true };
        image.outputFormat = { supported: false };
        image.maxOutputs = 1;
    } else if (protocol === "volcengine-ark-image") {
        image.references.maskSupported = false;
        image.quality.supported = false;
        image.transparentBackground.supported = false;
        image.responseFormat.supported = false;
        image.outputFormat.supported = false;
    }
    if (protocol === "volcengine-jimeng-image") {
        image.references.maxImages = 14;
        image.references.maskSupported = false;
        image.quality.supported = false;
        image.transparentBackground.supported = false;
        image.responseFormat.supported = false;
        image.outputFormat.supported = false;
    }
    if (protocol === "gemini-image") {
        image.references.maskSupported = false;
        // Gemini Images uses imageConfig.aspectRatio, not the OpenAI-style pixel size field.
        image.size = {
            parameter: "aspect_ratio",
            values: ["auto", "1:1", "2:3", "3:2", "3:4", "4:3", "9:16", "16:9", "21:9"],
            default: "1:1",
            allowCustom: false,
        };
        image.transparentBackground = { supported: false, default: false };
        image.responseFormat = { supported: false };
        image.outputFormat = { supported: false };
        image.maxOutputs = 4;
    }
    if (protocol === "dashscope-image") {
        // DashScope 支持比例和像素格式输入，适配器会自动转换
        image.size = {
            parameter: "size",
            values: [
                "auto",
                // 比例（推荐，更直观）
                "1:1", "3:2", "2:3", "4:3", "3:4", "16:9", "9:16",
                // 像素尺寸（高级用户）- 使用 x 格式与数据库存储一致
                "1024x1024", "1536x1024", "1024x1536", "2048x2048", "2048x1152",
            ],
            default: "1:1",
            allowCustom: false,
        };
        image.quality = { supported: true, values: ["auto", "low", "medium", "high"], default: "auto" };
    }
    if (protocol !== "grok-image" && model.trim().toLowerCase().startsWith("grok-imagine-image")) {
        image.references.maxImages = 0;
        image.references.maskSupported = false;
        image.size = {
            parameter: "aspect_ratio",
            values: ["1:1", "3:4", "4:3", "9:16", "16:9", "2:3", "3:2"],
            default: "1:1",
            allowCustom: false,
        };
        image.quality = { supported: true, values: ["1k", "2k"], default: "2k" };
        image.transparentBackground = { supported: false, default: false };
        image.responseFormat = { supported: true };
        image.outputFormat = { supported: false };
        image.maxOutputs = 1;
    }
    return image;
}

export function defaultModelCapabilityConfig(protocol?: ModelProtocol, model = ""): ModelCapabilityConfig {
    const text: TextCapabilityConfig = {
        // 文本模型的视觉能力必须由管理员明确开启，不能根据模型名猜测。
        references: { promptMaxChars: 32000, maxImages: 0, maxImageBytes: 0, maxVideos: 0, maxVideoBytes: 0 },
    };
    const video: VideoCapabilityConfig = {
        references: {
            promptMaxChars: 1000,
            minImages: 0,
            maxImages: 9,
            maxImageBytes: 30 * 1024 * 1024,
            maxVideos: 0,
            maxVideoBytes: 0,
            maxVideoDurationSeconds: 0,
            maxAudios: 0,
            maxAudioBytes: 0,
            maxAudioDurationSeconds: 0,
        },
        duration: { selection: "range", min: 1, max: 15, step: 1, default: 6 },
        ratios: ["16:9", "9:16", "1:1", "4:3", "3:4", "21:9"],
        defaultRatio: "16:9",
        resolutions: ["480p", "720p", "1080p", "1440p", "2160p"],
        defaultResolution: "720p",
        generateAudio: { supported: false, default: false },
        watermark: { supported: false, default: false },
        operations: ["text_to_video", "image_to_video"],
        defaultOperation: "text_to_video",
    };
    if (protocol === "volcengine-jimeng-video") {
        video.duration = { selection: "enum", values: [5, 10], default: 5 };
        video.resolutions = ["720p"];
    }
    if (protocol === "gemini-veo") {
        video.duration = { selection: "enum", values: [4, 6, 8], default: 6 };
        video.resolutions = ["720p", "1080p"];
    }
    if (protocol === "volcengine-ark-video" || protocol === "newapi-channel-1" || protocol === "newapi-channel-2") {
        video.references.maxVideos = 3;
        video.references.maxAudios = 3;
        video.references.maxVideoBytes = 200 * 1024 * 1024;
        video.references.maxAudioBytes = 15 * 1024 * 1024;
        video.references.maxVideoDurationSeconds = 15;
        video.references.maxAudioDurationSeconds = 15;
        video.generateAudio = { supported: true, default: true };
    }
    if (protocol === "volcengine-ark-video" || protocol === "newapi-channel-1") video.resolutions = ["480p", "720p", "1080p"];
    if (protocol === "volcengine-ark-video") video.watermark = { supported: true, default: false };
    if (protocol === "novita-video") {
        video.references.maxImages = 1;
        video.references.maxImageBytes = 10 * 1024 * 1024;
        video.duration = { selection: "enum", values: [5, 10], default: 5 };
        video.ratios = ["16:9", "9:16", "1:1"];
        video.resolutions = ["1080p"];
        video.defaultResolution = "1080p";
    }
    if (protocol === "dashscope-video") {
        // DashScope 视频协议（wan2.7, wan3.0, HappyHorse 系列）
        video.references.minImages = 0;
        video.references.maxImages = 9;
        video.references.maxImageBytes = 30 * 1024 * 1024;
        video.references.maxVideos = 3;
        video.references.maxVideoBytes = 50 * 1024 * 1024;
        video.references.maxVideoDurationSeconds = 30;
        video.duration = { selection: "range", min: 2, max: 16, step: 1, default: 6 };
        video.ratios = ["16:9", "9:16", "1:1", "4:3", "3:4"];
        video.defaultRatio = "16:9";
        video.resolutions = ["480p", "720p", "1080p"];
        video.defaultResolution = "720p";
        video.generateAudio = { supported: true, default: true };
        video.watermark = { supported: true, default: false };
        video.operations = ["text_to_video", "image_to_video", "reference_to_video"];
        video.defaultOperation = "text_to_video";
    }
    if (protocol === "dashscope-video" && model.trim().replace(/^models\//i, "").toLowerCase() === "wan3.0-video") {
        video.references = {
            promptMaxChars: 20000,
            minImages: 0,
            maxImages: 9,
            maxImageBytes: 30 * 1024 * 1024,
            maxVideos: 3,
            maxVideoBytes: 50 * 1024 * 1024,
            maxVideoDurationSeconds: 30,
            maxAudios: 1,
            maxAudioBytes: 15 * 1024 * 1024,
            maxAudioDurationSeconds: 30,
        };
        video.duration = { selection: "range", min: 2, max: 30, step: 1, default: 5, smartSupported: true };
        video.ratios = ["adaptive", "16:9", "4:3", "1:1", "3:4", "9:16"];
        video.defaultRatio = "adaptive";
        video.resolutions = ["480p", "720p", "1080p"];
        video.defaultResolution = "1080p";
        video.generateAudio = { supported: true, default: true };
        video.watermark = { supported: true, default: false };
        video.operations = ["text_to_video", "image_to_video", "reference_to_video", "audio_to_video"];
        video.defaultOperation = "text_to_video";
    }
    if (protocol === "dashscope-video") {
        const modelName = model.trim().replace(/^models\//i, "").toLowerCase();
        const wan27Kind = wan27VideoKind(modelName);
        if (wan27Kind) applyWan27VideoCapability(video, wan27Kind);
        const happyHorseKind = happyHorse11VideoKind(modelName);
        if (happyHorseKind) applyHappyHorseVideoCapability(video, happyHorseKind);
    }
    if (protocol === "minimax-video") {
        video.references.maxImages = 9;
        video.references.maxImageBytes = 30 * 1024 * 1024;
        video.references.maxVideos = 3;
        video.references.maxVideoBytes = 50 * 1024 * 1024;
        video.references.maxVideoDurationSeconds = 15;
        video.references.maxAudios = 3;
        video.references.maxAudioBytes = 15 * 1024 * 1024;
        video.references.maxAudioDurationSeconds = 15;
        video.duration = { selection: "enum", values: [4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15], default: 5 };
        video.ratios = ["adaptive", "21:9", "16:9", "4:3", "1:1", "3:4", "9:16"];
        video.resolutions = ["768P", "2K"];
        video.defaultResolution = "768P";
        video.watermark = { supported: true, default: false };
        video.operations.push("reference_to_video");
    }
    return { version: 1, text, image: defaultImageCapabilityConfig(protocol, model), video };
}

function wan27VideoKind(model: string) {
    if (["wan2.7-t2v", "wan2.7-t2v-2026-06-12", "wan2.7-t2v-2026-04-25"].includes(model)) return "t2v";
    if (["wan2.7-i2v", "wan2.7-i2v-2026-04-25"].includes(model)) return "i2v";
    if (["wan2.7-r2v", "wan2.7-r2v-2026-06-12"].includes(model)) return "r2v";
    return "";
}

function applyWan27VideoCapability(video: VideoCapabilityConfig, kind: string) {
    video.references = { promptMaxChars: 5000, minImages: 0, maxImages: 0, maxImageBytes: 0, maxVideos: 0, maxVideoBytes: 0, maxVideoDurationSeconds: 0, maxAudios: 0, maxAudioBytes: 0, maxAudioDurationSeconds: 0 };
    video.duration = { selection: "range", min: 2, max: 15, step: 1, default: 5 };
    video.ratios = ["16:9", "9:16", "1:1", "4:3", "3:4"];
    video.defaultRatio = "16:9";
    video.resolutions = ["720p", "1080p"];
    video.defaultResolution = "1080p";
    video.generateAudio = { supported: false, default: false };
    video.watermark = { supported: true, default: false };
    if (kind === "t2v") {
        Object.assign(video.references, { maxAudios: 1, maxAudioBytes: 15 * 1024 * 1024, maxAudioDurationSeconds: 30, audioMustFitOutput: true });
        video.operations = ["text_to_video", "audio_to_video"];
        video.defaultOperation = "text_to_video";
    } else if (kind === "i2v") {
        Object.assign(video.references, { minVisualReferences: 1, maxVisualReferences: 2, maxImages: 2, maxImageBytes: 20 * 1024 * 1024, maxVideos: 1, maxVideoBytes: 100 * 1024 * 1024, maxVideoDurationSeconds: 10, maxAudios: 1, maxAudioBytes: 15 * 1024 * 1024, maxAudioDurationSeconds: 30, audioMustFitOutput: true });
        video.ratios = ["adaptive"];
        video.defaultRatio = "adaptive";
        video.operations = ["image_to_video", "extend", "reference_to_video", "audio_to_video"];
        video.defaultOperation = "image_to_video";
    } else {
        Object.assign(video.references, { minVisualReferences: 1, maxVisualReferences: 5, maxImages: 5, maxImageBytes: 20 * 1024 * 1024, maxVideos: 5, maxVideoBytes: 100 * 1024 * 1024, maxVideoDurationSeconds: 30, maxOutputDurationWithVideoSeconds: 10, maxAudios: 1, maxAudioBytes: 15 * 1024 * 1024, maxAudioDurationSeconds: 10 });
        video.operations = ["reference_to_video", "image_to_video", "audio_to_video"];
        video.defaultOperation = "reference_to_video";
    }
}

function happyHorse11VideoKind(model: string) {
    if (model === "happyhorse-1.1-t2v") return "t2v";
    if (model === "happyhorse-1.1-i2v") return "i2v";
    if (model === "happyhorse-1.1-r2v") return "r2v";
    return "";
}

function applyHappyHorseVideoCapability(video: VideoCapabilityConfig, kind: string) {
    video.references = { promptMaxChars: 5000, minImages: 0, maxImages: 0, maxImageBytes: 0, maxVideos: 0, maxVideoBytes: 0, maxVideoDurationSeconds: 0, maxAudios: 0, maxAudioBytes: 0, maxAudioDurationSeconds: 0 };
    video.duration = { selection: "range", min: 3, max: 15, step: 1, default: 5 };
    video.ratios = ["16:9", "9:16", "1:1", "4:3", "3:4", "4:5", "5:4", "9:21", "21:9"];
    video.defaultRatio = "16:9";
    video.resolutions = ["480p", "720p", "1080p"];
    video.defaultResolution = "1080p";
    video.generateAudio = { supported: false, default: false };
    video.watermark = { supported: true, default: true };
    if (kind === "t2v") {
        video.operations = ["text_to_video"];
        video.defaultOperation = "text_to_video";
    } else if (kind === "i2v") {
        Object.assign(video.references, { minVisualReferences: 1, maxVisualReferences: 1, minImages: 1, maxImages: 1, maxImageBytes: 20 * 1024 * 1024 });
        video.ratios = ["adaptive"];
        video.defaultRatio = "adaptive";
        video.operations = ["image_to_video"];
        video.defaultOperation = "image_to_video";
    } else {
        Object.assign(video.references, { minVisualReferences: 1, maxVisualReferences: 9, minImages: 1, maxImages: 9, maxImageBytes: 20 * 1024 * 1024 });
        video.operations = ["reference_to_video", "image_to_video"];
        video.defaultOperation = "reference_to_video";
    }
}

export function modelCapabilityConfigFor(config: { channels: Array<{ id: string; models: string[]; modelCosts?: Array<{ model: string; capabilityConfig?: ModelCapabilityConfig; protocol?: ModelProtocol }> }> }, model: string) {
    const separator = model.indexOf("::");
    const channelId = separator >= 0 ? model.slice(0, separator) : "";
    const modelName = separator >= 0 ? model.slice(separator + 2) : model;
    const channel = config.channels.find((item) => item.id === channelId) || config.channels.find((item) => item.models.includes(modelName));
    const cost = channel?.modelCosts?.find((item) => item.model === modelName);
    const fallback = defaultModelCapabilityConfig(cost?.protocol, modelName);
    if (!cost?.capabilityConfig) return fallback;
    const capabilityConfig = normalizeModelCapabilityConfig(cost.capabilityConfig);
    const text = capabilityConfig.text ? { ...fallback.text!, ...capabilityConfig.text, references: { ...fallback.text!.references, ...capabilityConfig.text.references } } : fallback.text;
    const video = capabilityConfig.video ? { ...fallback.video!, ...capabilityConfig.video, references: { ...fallback.video!.references, ...capabilityConfig.video.references } } : fallback.video;
    const configuredImage = capabilityConfig.image;

    const image = configuredImage
        ? (() => {
              const configuredSize = configuredImage.size;
              const configuredValues = configuredSize?.values?.map(normalizeCapabilityString);
              const allowCustom = Boolean(configuredSize?.allowCustom || configuredValues?.includes("*"));
              const concreteValues = configuredValues?.filter((value) => value !== "*") || [];

              // Missing values inherit the protocol default. An explicit empty list
              // remains empty, and a wildcard enables custom input without inventing presets.
              const values = configuredValues === undefined ? fallback.image!.size.values : concreteValues;
              const configuredDefault = configuredSize?.default ? normalizeCapabilityString(configuredSize.default) : undefined;
              const defaultValue = configuredDefault && configuredDefault !== "*" && values.includes(configuredDefault) ? configuredDefault : values.find((value) => value !== "*") || fallback.image!.size.default;

              return {
                  ...fallback.image!,
                  ...configuredImage,
                  size: {
                      ...fallback.image!.size,
                      ...configuredSize,
                      values,
                      default: defaultValue,
                      allowCustom,
                  },
              };
          })()
        : fallback.image;
    return { ...fallback, ...capabilityConfig, text, image, video };
}

export function normalizeImageValue(profile: ImageCapabilityConfig, value: { size?: string; quality?: string; count?: string; transparentBackground?: string }) {
    const size = normalizeImageSizeSetting(profile, value.size);
    const quality = profile.quality.supported ? (value.quality && profile.quality.values.includes(value.quality) ? value.quality : profile.quality.default || "auto") : profile.quality.default || "auto";
    const count = String(Math.max(1, Math.min(profile.maxOutputs, Math.floor(Math.abs(Number(value.count)) || 1))));
    const transparentBackground = profile.transparentBackground.supported && value.transparentBackground === "true" ? "true" : "false";
    return { size, quality, count, transparentBackground };
}

export function normalizeImageSizeSetting(profile: ImageCapabilityConfig, value?: string) {
    if (profile.size.parameter === "none") return "auto";
    if (!profile.size.allowCustom && profile.size.values.length === 0) return "auto";
    const candidate = value?.trim() || profile.size.default;
    if (profile.size.allowCustom || profile.size.values.includes(candidate)) return candidate;
    return profile.size.default || profile.size.values[0] || "auto";
}

export function imageSizeRequest(profile: ImageCapabilityConfig, value?: string) {
    const parameter = profile.size.parameter;
    if (parameter === "none") return undefined;
    const normalized = normalizeImageSizeSetting(profile, value);
    if (!normalized || normalized === "auto") return undefined;
    return { parameter, value: normalized };
}

export function normalizeVideoValue(profile: VideoCapabilityConfig, value: { seconds?: string; ratio?: string; resolution?: string }) {
    const requestedDuration = Number(value.seconds);
    const duration = requestedDuration === -1 && profile.duration.smartSupported
        ? -1
        : profile.duration.selection === "enum"
          ? ((profile.duration.values || []).includes(requestedDuration) ? requestedDuration : profile.duration.default)
          : normalizeRangeDuration(profile, requestedDuration);
    const ratio = profile.ratios.includes(value.ratio || "") ? value.ratio! : profile.defaultRatio;
    // 前端状态历史上保存过 `720`，而能力配置和供应商通常使用 `720p`；统一按能力中的原始值返回，避免被误判为不支持。
    const resolution = videoResolutionRequest(profile, value.resolution) || profile.defaultResolution || profile.resolutions[0] || "";
    return { seconds: String(duration), ratio, resolution };
}

export function videoResolutionRequest(profile: VideoCapabilityConfig, value: string | undefined) {
    const requested = String(value || "")
        .trim()
        .toLowerCase();
    if (!requested || requested === "auto" || requested === "default" || requested === "medium" || requested === "high") return undefined;
    const candidates = [requested];
    if (/^\d+$/.test(requested)) candidates.push(`${requested}p`);
    if (requested === "low") candidates.push("480p");
    if (requested === "2k") candidates.push("1440p");
    if (requested === "1440" || requested === "1440p") candidates.push("2k");
    if (requested === "4k") candidates.push("2160p");
    if (requested === "2160" || requested === "2160p") candidates.push("4k");
    const supported = new Map(profile.resolutions.map((resolution) => [resolution.trim().toLowerCase(), resolution.trim()]));
    for (const candidate of candidates) {
        const match = supported.get(candidate);
        if (match) return match;
    }
    return undefined;
}

function normalizeRangeDuration(profile: VideoCapabilityConfig, value: number) {
    const min = profile.duration.min || 1;
    const max = profile.duration.max || min;
    const step = profile.duration.step || 1;
    const candidate = Number.isFinite(value) ? Math.floor(value) : profile.duration.default;
    const clamped = Math.min(max, Math.max(min, candidate));
    const maxStep = Math.max(0, Math.floor((max - min) / step));
    return min + Math.min(maxStep, Math.max(0, Math.round((clamped - min) / step))) * step;
}

export function videoDurationOptions(profile: VideoCapabilityConfig) {
    if (profile.duration.selection === "enum") return profile.duration.values || [];
    const min = profile.duration.min || 1;
    const max = profile.duration.max || min;
    const step = profile.duration.step || 1;
    return Array.from({ length: Math.floor((max - min) / step) + 1 }, (_, index) => min + index * step);
}

export function videoDurationAllowed(profile: VideoCapabilityConfig, value: number) {
    if (value === -1) return profile.duration.smartSupported === true;
    if (profile.duration.selection === "enum") return (profile.duration.values || []).includes(value);
    const min = profile.duration.min || 1;
    const max = profile.duration.max || min;
    const step = profile.duration.step || 1;
    return value >= min && value <= max && (value - min) % step === 0;
}
