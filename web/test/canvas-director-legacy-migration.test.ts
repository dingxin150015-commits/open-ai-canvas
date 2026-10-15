import { expect, test } from "bun:test";

import { normalizeLegacyDirectorProject } from "../src/lib/canvas/canvas-director-legacy-migration";
import type { CanvasProject } from "../src/stores/canvas/use-canvas-store";

test("retains old Director scene data while exposing it to Previs", () => {
    const scene = { id: "scene-old", objects: [], cameras: [], lights: [], shots: [{ id: "shot-old" }] };
    const source = {
        id: "canvas-old",
        directorScenes: [scene],
        nodes: [{ id: "node-old", metadata: { directorSceneId: "scene-old", directorShotId: "shot-old", directorPreviewNodeId: "preview-old" } }],
    } as unknown as CanvasProject;

    const migrated = normalizeLegacyDirectorProject(source);
    expect(migrated.previsScenes).toHaveLength(1);
    expect(migrated.previsScenes[0].id).toBe("scene-old");
    expect(migrated.nodes[0].metadata).toMatchObject({ previsSceneId: "scene-old", previsShotId: "shot-old", previsPreviewNodeId: "preview-old" });
    expect((migrated as CanvasProject & { directorScenes?: unknown }).directorScenes).toEqual([scene]);
    expect(source.nodes[0].metadata).not.toHaveProperty("previsSceneId");
    expect(normalizeLegacyDirectorProject(migrated).previsScenes).toHaveLength(1);
});

test("does not reinterpret malformed legacy scenes or overwrite a newer Previs scene", () => {
    const scene = { id: "scene-old", objects: [], cameras: [], lights: [], shots: [] };
    const current = { ...scene, title: "newer" };
    const source = { directorScenes: [scene], previsScenes: [current], nodes: [] } as unknown as CanvasProject;
    const migrated = normalizeLegacyDirectorProject(source);
    expect(migrated.previsScenes).toEqual([current]);

    const malformed = { directorScenes: [{ id: "incomplete" }], nodes: [] } as unknown as CanvasProject;
    expect(normalizeLegacyDirectorProject(malformed)).toBe(malformed);
});
