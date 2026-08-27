import { describe, expect, test } from "bun:test";

import { defaultImageCapabilityConfig, imageSizeRequest, normalizeImageValue } from "../src/lib/model-capabilities";

describe("Qwen Image 3.0 capability", () => {
    test("uses the official conservative Qwen 3.0 profile", () => {
        const profile = defaultImageCapabilityConfig("dashscope-image", "qwen-image-3.0-pro");

        expect(profile.references).toEqual({
            promptMaxChars: 32000,
            maxImages: 3,
            maxImageBytes: 10 * 1024 * 1024,
            maskSupported: false,
        });
        expect(profile.size.default).toBe("auto");
        expect(profile.size.allowCustom).toBe(true);
        expect(profile.size.values).toContain("1024x1536");
        expect(profile.quality.supported).toBe(false);
        expect(profile.transparentBackground.supported).toBe(false);
        expect(profile.responseFormat.supported).toBe(false);
        expect(profile.outputFormat.supported).toBe(false);
        expect(profile.maxOutputs).toBe(6);
    });

    test("treats auto as size omission and preserves supported custom pixels and ratios", () => {
        const profile = defaultImageCapabilityConfig("dashscope-image", "qwen-image-3.0-pro");

        expect(imageSizeRequest(profile, "auto")).toBeUndefined();
        expect(imageSizeRequest(profile, "1024x1536")).toEqual({ parameter: "size", value: "1024x1536" });
        expect(imageSizeRequest(profile, "5:4")).toEqual({ parameter: "size", value: "5:4" });
        expect(normalizeImageValue(profile, { size: "auto", quality: "high", count: "9", transparentBackground: "true" })).toEqual({
            size: "auto",
            quality: "auto",
            count: "6",
            transparentBackground: "false",
        });
    });

    test("does not leak the Qwen contract to another DashScope image model", () => {
        const profile = defaultImageCapabilityConfig("dashscope-image", "wan2.7-image-pro");

        expect(profile.references.maxImages).not.toBe(3);
        expect(profile.maxOutputs).not.toBe(6);
    });
});
