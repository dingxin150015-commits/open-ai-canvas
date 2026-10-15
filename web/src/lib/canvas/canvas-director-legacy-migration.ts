import type { CanvasProject } from "@/stores/canvas/use-canvas-store";
import type { PrevisScene } from "@/types/previs";

type LegacyProject = CanvasProject & { directorScenes?: unknown };

const metadataKeys = [
    ["directorSceneId", "previsSceneId"],
    ["directorShotId", "previsShotId"],
    ["directorPreviewNodeId", "previsPreviewNodeId"],
    ["directorDepthNodeId", "previsDepthNodeId"],
    ["directorNormalNodeId", "previsNormalNodeId"],
    ["directorClayVideoNodeId", "previsClayVideoNodeId"],
] as const;

function isLegacyScene(value: unknown): value is PrevisScene {
    if (!value || typeof value !== "object") return false;
    const scene = value as Partial<PrevisScene>;
    return typeof scene.id === "string" && scene.id.length > 0 && Array.isArray(scene.objects) && Array.isArray(scene.cameras) && Array.isArray(scene.lights) && Array.isArray(scene.shots);
}

// The old Director scene uses the same scene geometry as Previs, but its
// project and node metadata keys changed. Keep the original fields in the
// serialized project so a failed migration never destroys the source data.
export function normalizeLegacyDirectorProject(project: CanvasProject): CanvasProject {
    const legacy = (project as LegacyProject).directorScenes;
    if (!Array.isArray(legacy) || !legacy.length || !legacy.every(isLegacyScene) || !Array.isArray(project.nodes)) return project;

    const scenes = Array.isArray(project.previsScenes) ? [...project.previsScenes] : [];
    const sceneIds = new Set(scenes.map((scene) => scene.id));
    for (const scene of legacy) {
        if (sceneIds.has(scene.id)) continue;
        scenes.push(scene);
        sceneIds.add(scene.id);
    }

    const nodes = project.nodes.map((node) => {
        const metadata = node.metadata as Record<string, unknown> | undefined;
        if (!metadata || typeof metadata.directorSceneId !== "string" || !sceneIds.has(metadata.directorSceneId)) return node;
        const patch: Record<string, string> = {};
        for (const [oldKey, newKey] of metadataKeys) {
            const oldValue = metadata[oldKey];
            if (typeof metadata[newKey] !== "string" && typeof oldValue === "string") patch[newKey] = oldValue;
        }
        return Object.keys(patch).length ? { ...node, metadata: { ...node.metadata, ...patch } } : node;
    });

    return { ...project, nodes, previsScenes: scenes };
}
