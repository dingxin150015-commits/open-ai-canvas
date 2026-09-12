import { describe, expect, test } from "bun:test";
import { bailianReferencePrompt, optionalVideoPrompt, videoAttachmentLimit, videoAttachmentError, videoFrameRoles, validateBailianOptions } from "../src/lib/bailian-generation";
import { defaultModelCapabilityConfig } from "../src/lib/model-capabilities";
import { parseCreationSettingsDraftDocument } from "../src/services/creation-settings-store";

describe("Bailian official capability boundaries", () => {
    test.each([
        ["wan2.7-t2v", ["audio"], 1],
        ["wan2.7-i2v", ["image", "image", "audio"], 3],
        ["wan2.7-r2v", ["image", "image", "video", "video", "video", "audio"], 6],
        ["happyhorse-1.1-r2v", Array(9).fill("image"), 9],
    ])("%s permits the complete reference bundle", (model, kinds, limit) => {
        const profile = defaultModelCapabilityConfig("dashscope-video", model).video!;
        expect(videoAttachmentLimit(profile)).toBe(limit);
        expect(videoAttachmentError(profile, kinds)).toBe("");
        expect(videoAttachmentError(profile, [...kinds, "audio"])).not.toBe("");
    });
    test("optional prompts are exact-model and media dependent", () => {
        expect(optionalVideoPrompt("c::happyhorse-1.1-i2v", 1, 0)).toBe(true);
        expect(optionalVideoPrompt("wan2.7-i2v", 0, 1)).toBe(true);
        for (const model of ["wan2.7-t2v", "wan2.7-r2v", "other-i2v"]) expect(optionalVideoPrompt(model, 1, 0)).toBe(false);
        expect(optionalVideoPrompt("wan2.7-i2v", 0, 0)).toBe(false);
    });
    test("frame roles match each model", () => {
        expect(videoFrameRoles("happyhorse-1.1-r2v")).toEqual({ first: false, last: false });
        expect(videoFrameRoles("happyhorse-1.1-i2v")).toEqual({ first: true, last: false });
        expect(videoFrameRoles("wan2.7-r2v")).toEqual({ first: true, last: false });
        expect(videoFrameRoles("wan2.7-i2v")).toEqual({ first: true, last: true });
    });
    test("references are typed and first-frame is not a reference-image ordinal", () => {
        expect(bailianReferencePrompt("@视频1 和 @图片2，@图片1", "wan2.7-r2v", [{ id: "first" }, { id: "ref" }], "first")).toBe("Video 1 和 Image 1，首帧");
        expect(bailianReferencePrompt("@图片2 @图片1", "happyhorse-1.1-r2v", [{ id: "a" }, { id: "b" }])).toBe("[Image 2] [Image 1]");
    });
    test("advanced false values survive draft persistence and invalid seed fails", () => {
        const generationOptions = { seed: "0", prompt_extend: false, enable_thinking: false, watermark: false, request_mode: "async" };
        const parsed = parseCreationSettingsDraftDocument(JSON.stringify({ version: 1, drafts: { key: { ratio: "auto", generationOptions } } }));
        expect(parsed.drafts.key.generationOptions).toEqual(generationOptions);
        expect(() => validateBailianOptions({ seed: "1.5" }, "qwen-image-3.0-pro", 0)).toThrow();
        expect(() => validateBailianOptions({ prompt_extend_mode: "agent" }, "qwen-image-3.0-pro", 1)).toThrow();
    });
});
