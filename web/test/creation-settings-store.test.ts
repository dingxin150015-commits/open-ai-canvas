import assert from "node:assert/strict";
import { afterEach, test } from "node:test";

import localforage from "localforage";

import { CREATION_SETTINGS_DRAFTS_KEY, creationSettingsDraftKey, loadCreationSettingsDraft, parseCreationSettingsDraftDocument, saveCreationSettingsDraft } from "../src/services/creation-settings-store";

const originalWindow = globalThis.window;
const originalGetItem = localforage.getItem.bind(localforage);
const originalSetItem = localforage.setItem.bind(localforage);

afterEach(() => {
    globalThis.window = originalWindow;
    localforage.getItem = originalGetItem;
    localforage.setItem = originalSetItem;
});

test("creation settings draft identity separates model, mode, and operation", () => {
    assert.notEqual(creationSettingsDraftKey({ mode: "video", model: "channel::wan3.0-video", operation: "text_to_video" }), creationSettingsDraftKey({ mode: "video", model: "channel::wan3.0-video", operation: "image_to_video" }));
    assert.notEqual(creationSettingsDraftKey({ mode: "video", model: "channel::wan3.0-video", operation: "text_to_video" }), creationSettingsDraftKey({ mode: "video", model: "channel::wan2.7-t2v", operation: "text_to_video" }));
});

test("creation settings parser fails closed for malformed documents and bounds fields", () => {
    assert.deepEqual(parseCreationSettingsDraftDocument("not-json"), { version: 1, drafts: {} });
    const key = creationSettingsDraftKey({ mode: "video", model: "model", operation: "text_to_video" });
    const parsed = parseCreationSettingsDraftDocument(
        JSON.stringify({
            version: 1,
            drafts: {
                [key]: {
                    ratio: " 16:9 ",
                    seconds: "2",
                    videoQuality: "480",
                    videoGenerateAudio: false,
                    videoWatermark: true,
                    ignored: "secret",
                },
            },
        }),
    );
    assert.deepEqual(parsed.drafts[key], {
        ratio: "16:9",
        seconds: "2",
        videoQuality: "480",
        videoGenerateAudio: false,
        videoWatermark: true,
    });
});

test("creation settings writes serialize per user scope without losing another operation", async () => {
    const values = new Map<string, string>();
    globalThis.window = {} as Window & typeof globalThis;
    localforage.getItem = (async (key: string) => values.get(key) ?? null) as typeof localforage.getItem;
    localforage.setItem = (async (key: string, value: string) => {
        values.set(key, value);
        return value;
    }) as typeof localforage.setItem;

    const textToVideo = { mode: "video" as const, model: "channel::wan3.0-video", operation: "text_to_video" };
    const imageToVideo = { mode: "video" as const, model: "channel::wan3.0-video", operation: "image_to_video" };
    await Promise.all([
        saveCreationSettingsDraft(textToVideo, { ratio: "16:9", seconds: "2", videoQuality: "480", videoGenerateAudio: false, videoWatermark: false }, "user-a"),
        saveCreationSettingsDraft(imageToVideo, { ratio: "9:16", seconds: "5", videoQuality: "720", videoGenerateAudio: true, videoWatermark: true }, "user-a"),
    ]);

    assert.deepEqual(await loadCreationSettingsDraft(textToVideo, "user-a"), {
        ratio: "16:9",
        seconds: "2",
        videoQuality: "480",
        videoGenerateAudio: false,
        videoWatermark: false,
    });
    assert.deepEqual(await loadCreationSettingsDraft(imageToVideo, "user-a"), {
        ratio: "9:16",
        seconds: "5",
        videoQuality: "720",
        videoGenerateAudio: true,
        videoWatermark: true,
    });
    assert.equal(values.has(`${CREATION_SETTINGS_DRAFTS_KEY}:user:user-a`), true);
    assert.equal(values.has(`${CREATION_SETTINGS_DRAFTS_KEY}:user:user-b`), false);
});
