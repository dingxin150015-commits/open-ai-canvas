import { expect, test } from "bun:test";

import { generationArtifactStorageKey } from "../src/services/generation-artifact-sink";
import { localArtifactStoreForStorageKey } from "../src/services/local-artifact-storage-key";

type Scenario = "image-cleanup" | "scope-cleanup-switch" | "scope-cleanup-late-canvas-reference" | "video-commit-race" | "audio-commit-race";
type ScenarioResponse<T> = { ok: true; result: T } | { ok: false; error: string };

function runScenario<T>(scenario: Scenario): Promise<T> {
    return new Promise<T>((resolve, reject) => {
        const worker = new Worker(new URL("./helpers/generation-storage-consistency.worker.ts", import.meta.url).href, { type: "module" });
        const finish = () => {
            worker.terminate();
        };
        worker.onmessage = (event: MessageEvent<ScenarioResponse<T>>) => {
            finish();
            if (event.data.ok) resolve(event.data.result);
            else reject(new Error(event.data.error));
        };
        worker.onerror = (event) => {
            finish();
            reject(event.error ?? new Error(event.message));
        };
        worker.postMessage(scenario);
    });
}

test("generation image cleanup removes unused generation-image blobs and preserves referenced ones", async () => {
    const result = await runScenario<{ usedPresent: boolean; unusedPresent: boolean }>("image-cleanup");
    expect(result.usedPresent).toBe(true);
    expect(result.unusedPresent).toBe(false);
});

test("delayed cleanup keeps account-A Canvas-only image references after switching the active account", async () => {
    const result = await runScenario<{ referencedPresent: boolean; unusedPresent: boolean; otherScopePresent: boolean }>("scope-cleanup-switch");
    expect(result.referencedPresent).toBe(true);
    expect(result.unusedPresent).toBe(false);
    expect(result.otherScopePresent).toBe(true);
});

test("delayed cleanup preserves a same-scope Canvas-only reference added after cleanup was queued", async () => {
    const result = await runScenario<{ referencedPresent: boolean; unusedPresent: boolean }>("scope-cleanup-late-canvas-reference");
    expect(result.referencedPresent).toBe(true);
    expect(result.unusedPresent).toBe(false);
});

test("generation video and audio materialization cannot race cleanup into a catalog row whose blob is missing", async () => {
    for (const mediaType of ["video", "audio"] as const) {
        const result = await runScenario<{ kind?: string; storageKey?: string; blobPresent: boolean; generationAssetCount: number }>(`${mediaType}-commit-race`);
        expect(result.kind).toBe(mediaType);
        expect(result.storageKey).toMatch(new RegExp(`^generation-${mediaType}:generation-media-commit-race-${mediaType}:`));
        expect(result.blobPresent).toBe(true);
        expect(result.generationAssetCount).toBe(1);
    }
}, 15_000);

test("remote sync routes generated artifacts to the store that materialized them", () => {
    expect(localArtifactStoreForStorageKey(generationArtifactStorageKey("effect-image", "image", "user-1"))).toBe("image");
    expect(localArtifactStoreForStorageKey(generationArtifactStorageKey("effect-video", "video", "user-1"))).toBe("media");
    expect(localArtifactStoreForStorageKey(generationArtifactStorageKey("effect-audio", "audio", "user-1"))).toBe("media");
});

test("remote sync preserves existing local artifact key families", () => {
    for (const key of ["image:user-1:id", "generation-image:user-1:id"]) expect(localArtifactStoreForStorageKey(key)).toBe("image");
    for (const key of ["video:user-1:id", "audio:user-1:id", "file:user-1:id", "video-reference:user-1:id", "audio-reference:user-1:id", "generation-video:user-1:id", "generation-audio:user-1:id"])
        expect(localArtifactStoreForStorageKey(key)).toBe("media");
    for (const key of ["", "resource:remote-id", "blob:https://example.com/id", "https://example.com/video.mp4"]) expect(localArtifactStoreForStorageKey(key)).toBeNull();
});
