import { describe, expect, test } from "bun:test";

import { defaultModelCapabilityConfig } from "../src/lib/model-capabilities";

describe("DashScope model-specific video capabilities", () => {
    test("Wan 2.7 T2V/I2V/R2V profiles remain distinct", () => {
        const t2v = defaultModelCapabilityConfig("dashscope-video", "wan2.7-t2v-2026-06-12").video!;
        const i2v = defaultModelCapabilityConfig("dashscope-video", "wan2.7-i2v-2026-04-25").video!;
        const r2v = defaultModelCapabilityConfig("dashscope-video", "wan2.7-r2v").video!;
        expect(t2v.references.maxImages).toBe(0);
        expect(t2v.references.maxAudios).toBe(1);
        expect(t2v.references.audioMustFitOutput).toBe(true);
        expect(t2v.resolutions).toEqual(["720p", "1080p"]);
        expect(i2v.references.maxImages).toBe(2);
        expect(i2v.references.maxVideos).toBe(1);
        expect(i2v.references.audioMustFitOutput).toBe(true);
        expect(i2v.references.maxVisualReferences).toBe(2);
        expect(i2v.ratios).toEqual(["adaptive"]);
        expect(i2v.operations).toContain("extend");
        expect(r2v.references.maxImages).toBe(5);
        expect(r2v.references.maxVideos).toBe(5);
        expect(r2v.references.maxAudioDurationSeconds).toBe(10);
        expect(r2v.references.maxOutputDurationWithVideoSeconds).toBe(10);
        expect(r2v.references.maxVisualReferences).toBe(5);
    });

    test("HappyHorse 1.1 profiles use exact media contracts", () => {
        const t2v = defaultModelCapabilityConfig("dashscope-video", "happyhorse-1.1-t2v").video!;
        const i2v = defaultModelCapabilityConfig("dashscope-video", "happyhorse-1.1-i2v").video!;
        const r2v = defaultModelCapabilityConfig("dashscope-video", "happyhorse-1.1-r2v").video!;
        expect(t2v.duration.min).toBe(3);
        expect(t2v.watermark).toEqual({ supported: true, default: true });
        expect(t2v.ratios).toHaveLength(9);
        expect(i2v.references).toMatchObject({ minImages: 1, maxImages: 1, maxVideos: 0, maxAudios: 0 });
        expect(i2v.ratios).toEqual(["adaptive"]);
        expect(r2v.references).toMatchObject({ minImages: 1, maxImages: 9, maxVideos: 0, maxAudios: 0 });
    });

    test("HappyHorse 1.0 remains on the non-ready fallback", () => {
        const profile = defaultModelCapabilityConfig("dashscope-video", "happyhorse-1.0-t2v").video!;
        expect(profile.duration.min).toBe(2);
        expect(profile.duration.max).toBe(16);
        expect(profile.watermark.default).toBe(false);
    });
});
