import { modelCapabilityConfigFor, normalizeVideoValue } from "@/lib/model-capabilities";
import { resolveModelVideoBooleanOptions } from "@/lib/model-selection";
import type { AiConfig } from "@/stores/use-config-store";
import type { CreationSettings } from "./creation-types";

// 报价与提交共用所选模型的音频、水印和尺寸约束。
export function creationVideoConfig(config: AiConfig, model: string, settings: Pick<CreationSettings, "ratio" | "seconds" | "videoQuality"> & Partial<Pick<CreationSettings, "videoGenerateAudio" | "videoWatermark">>): AiConfig {
    const profile = modelCapabilityConfigFor(config, model).video!;
    const normalized = normalizeVideoValue(profile, {
        ratio: settings.ratio,
        seconds: settings.seconds,
        resolution: settings.videoQuality,
    });
    return {
        ...config,
        model,
        videoModel: model,
        size: profile.ratios.length > 0 ? normalized.ratio : "",
        videoSeconds: normalized.seconds,
        vquality: normalized.resolution.replace(/p$/i, ""),
        ...resolveModelVideoBooleanOptions(
            config,
            model,
            {
                ...(settings.videoGenerateAudio === undefined ? {} : { videoGenerateAudio: String(settings.videoGenerateAudio) }),
                ...(settings.videoWatermark === undefined ? {} : { videoWatermark: String(settings.videoWatermark) }),
            },
            {
                videoGenerateAudio: config.videoGenerateAudio,
                videoWatermark: config.videoWatermark,
            },
        ),
    };
}
