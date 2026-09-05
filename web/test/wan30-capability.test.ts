import { describe, expect, test } from "bun:test";

import { defaultModelCapabilityConfig, normalizeVideoValue, videoDurationAllowed } from "../src/lib/model-capabilities";
import { inferVideoOperation } from "../src/lib/model-selection";

describe("Wan 3.0 video capability", () => {
    test("uses the official model-specific contract", () => {
        const profile = defaultModelCapabilityConfig("dashscope-video", "wan3.0-video").video!;
        expect(profile.duration).toEqual({ selection: "range", min: 2, max: 30, step: 1, default: 5, smartSupported: true });
        expect(profile.ratios).toEqual(["adaptive", "16:9", "4:3", "1:1", "3:4", "9:16"]);
        expect(profile.resolutions).toEqual(["480p", "720p", "1080p"]);
        expect(profile.defaultResolution).toBe("1080p");
        expect(profile.generateAudio).toEqual({ supported: true, default: true });
        expect(profile.operations).toContain("audio_to_video");
    });

    test("preserves smart duration and rejects unsupported duration", () => {
        const profile = defaultModelCapabilityConfig("dashscope-video", "wan3.0-video").video!;
        expect(normalizeVideoValue(profile, { seconds: "-1", ratio: "adaptive", resolution: "1080p" }).seconds).toBe("-1");
        expect(videoDurationAllowed(profile, -1)).toBe(true);
        expect(videoDurationAllowed(profile, 2)).toBe(true);
        expect(videoDurationAllowed(profile, 30)).toBe(true);
        expect(videoDurationAllowed(profile, 31)).toBe(false);
    });

    test("does not apply the Wan 3.0 profile to Prime while it remains planned", () => {
        const profile = defaultModelCapabilityConfig("dashscope-video", "wan3.0-video-prime").video!;
        expect(profile.duration.max).toBe(16);
        expect(profile.duration.smartSupported).not.toBe(true);
    });

    test("routes mixed references to reference mode", () => {
        expect(inferVideoOperation({ textCount: 0, imageCount: 1, videoCount: 1, audioCount: 1, characterCount: 0 })).toBe("reference_to_video");
        expect(inferVideoOperation({ textCount: 0, imageCount: 1, videoCount: 1, audioCount: 0, characterCount: 0 })).toBe("reference_to_video");
        expect(inferVideoOperation({ textCount: 0, imageCount: 1, videoCount: 0, audioCount: 0, characterCount: 0 })).toBe("image_to_video");
    });
});
