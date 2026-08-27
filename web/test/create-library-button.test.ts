import { describe, expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { resolve } from "node:path";

function compactSource(source: string) {
    return source.replace(/\s+/g, " ").trim();
}

describe("creation library button", () => {
    test("places a library control beside the generation mode picker", () => {
        const source = readFileSync(resolve(import.meta.dir, "../src/pages/create/index.tsx"), "utf8");
        const dockStart = source.indexOf('<footer className="creation-chat-dock">');
        const dockEnd = source.indexOf("</footer>", dockStart);

        expect(dockStart).toBeGreaterThanOrEqual(0);
        expect(dockEnd).toBeGreaterThan(dockStart);
        const dockSource = compactSource(source.slice(dockStart, dockEnd));
        const modePickerIndex = dockSource.indexOf("<ModePicker mode={props.mode}");
        const attachmentIndex = dockSource.indexOf('aria-label="从本机上传附件"');
        const libraryIndex = dockSource.indexOf('aria-label="打开素材库选择参考内容"');

        expect(modePickerIndex).toBeGreaterThanOrEqual(0);
        expect(attachmentIndex).toBeGreaterThan(modePickerIndex);
        expect(libraryIndex).toBeGreaterThan(attachmentIndex);
        expect(dockSource).toContain("onClick={props.onOpenLibrary}");
        expect(dockSource).toContain("disabled={props.busy || !referencesSupported}");
    });

    test("uploads from the library without adding a reference before confirmation", () => {
        const source = readFileSync(resolve(import.meta.dir, "../src/pages/create/index.tsx"), "utf8");
        const pickerSource = readFileSync(resolve(import.meta.dir, "../src/components/assets/asset-library-picker-modal.tsx"), "utf8");
        const uploadStart = source.indexOf("const uploadLibraryAssets = async");
        const uploadEnd = source.indexOf("const handleFileChange", uploadStart);

        expect(uploadStart).toBeGreaterThanOrEqual(0);
        expect(uploadEnd).toBeGreaterThan(uploadStart);
        expect(source.slice(uploadStart, uploadEnd)).not.toContain("setAttachments");
        expect(source).toContain("onUpload: uploadLibraryAssets");
        expect(source).not.toContain("onUpload={() => fileInputRef.current?.click()}");
        expect(source).toContain("上传后保存到素材库");
        expect(pickerSource).toContain("保存完成后会自动选中");
        expect(source).toContain("个素材已上传到素材库并自动选中");
    });

    test("previews prompt reference images without removing them", () => {
        const createSource = readFileSync(resolve(import.meta.dir, "../src/pages/create/index.tsx"), "utf8");
        const canvasSource = readFileSync(resolve(import.meta.dir, "../src/components/canvas/canvas-node-prompt-panel.tsx"), "utf8");

        expect(createSource).toContain('className="creation-user-message-attachments"');
        expect(createSource).toContain('setPreviewType(kind === "video" ? "video" : "image")');
        expect(createSource).toContain("<CreationMediaPreviewModal url={previewUrl} type={previewType}");
        expect(canvasSource).toContain("canPreview ? setImagePreview(reference) : onInsert(reference)");
        expect(canvasSource).toContain("<AntImage");
        expect(canvasSource).toContain("onClick={() => onInsert(reference)}");
    });

    test("promotes the first reference into the primary reference slot", () => {
        const source = readFileSync(resolve(import.meta.dir, "../src/pages/create/index.tsx"), "utf8");

        expect(source).toContain("const [primaryAttachment, ...secondaryAttachments] = props.attachments");
        expect(source).toContain("<CreationAttachmentThumbnail item={primaryAttachment} primary");
        expect(source).toContain("secondaryAttachments.map((item) => <CreationAttachmentThumbnail");
        expect(source).toContain('className={primary ? "creation-chat-reference is-paper creation-chat-reference-media" : "creation-chat-attachment"}');
    });

    test("shows capability-driven audio and watermark controls in the real Create settings menu", () => {
        const source = readFileSync(resolve(import.meta.dir, "../src/pages/create/index.tsx"), "utf8");
        const menuStart = source.indexOf("function GenerationSettingsMenu");
        const menuEnd = source.indexOf("function SettingSection", menuStart);
        const menuSource = compactSource(source.slice(menuStart, menuEnd));

        expect(menuStart).toBeGreaterThanOrEqual(0);
        expect(menuEnd).toBeGreaterThan(menuStart);
        expect(menuSource).toContain("props.videoProfile.generateAudio.supported");
        expect(menuSource).toContain("props.videoProfile.watermark.supported");
        expect(menuSource).toContain('onChange={props.onVideoGenerateAudioChange} aria-label="生成声音"');
        expect(menuSource).toContain('onChange={props.onVideoWatermarkChange} aria-label="添加水印"');
        expect(menuSource).toContain('className="creation-output-switches"');
        expect(menuSource).toContain("生成声音");
        expect(menuSource).toContain("添加水印");
        expect(menuSource).toContain('videoGenerateAudio ? "有声" : "无声"');
        expect(menuSource).toContain('videoWatermark ? "有水印" : "无水印"');
    });

    test("persists Create settings by user, model, and operation and submits the visible output values", () => {
        const source = compactSource(readFileSync(resolve(import.meta.dir, "../src/pages/create/index.tsx"), "utf8"));

        expect(source).toContain("loadCreationSettingsDraft(settingsDraftIdentity, creationSettingsScope)");
        expect(source).toContain("saveCreationSettingsDraft(settingsDraftIdentity, draft, creationSettingsScope)");
        expect(source).toContain("operation: settingsOperation");
        expect(source).toContain("videoGenerateAudio: String(videoProfile.generateAudio.supported && videoGenerateAudio)");
        expect(source).toContain("videoWatermark: String(videoProfile.watermark.supported && videoWatermark)");
        expect(source).toContain("const canSubmit = Boolean(props.prompt.trim()) && !props.busy && props.settingsReady");
    });
});
