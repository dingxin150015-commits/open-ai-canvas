import { describe, expect, test } from "bun:test";
import { executeImageGeneration } from "../src/pages/canvas/canvas-image-generation-executor";
import type { CanvasGenerationExecution } from "../src/pages/canvas/canvas-generation-executor-types";
import { createModelChannel, defaultConfig, type AiConfig } from "../src/stores/use-config-store";
import { nativeImageBatchEnabled } from "../src/lib/bailian-generation";
import { applyGenerationTaskResultToNodes, applyMaterializedGenerationTaskResultToNodes } from "../src/lib/canvas/canvas-generation-task-sync";
import { buildNodeConfig } from "../src/components/canvas/canvas-node-prompt-panel";
import { useAssetStore, type ImageAsset } from "../src/stores/use-asset-store";
import { createGenerationRetryContext } from "../src/lib/canvas/canvas-project-generation";
import { CanvasNodeType, type CanvasNodeData } from "../src/types/canvas";
import type { GenerationTask } from "../src/services/api/task-center";

const config: AiConfig = { ...defaultConfig, model: "qwen-image-3.0-pro", count: "3", generationOptions: { contractKind: "qwen", nativeBatch: true } };
const task: GenerationTask = {
    id: "shared-task", type: "canvas_image", status: "succeeded", prompt: "fixture", attempts: 1, createdAt: "", updatedAt: "",
    resultJson: JSON.stringify({ mode: "image", images: [0, 1, 2].map((index) => ({ dataUrl: `https://fixture.invalid/${index}`, storageKey: `resource:image-${index}`, width: 1024, height: 1024 })) }),
};

async function execute(overrides: Partial<AiConfig> = {}, failure?: "submit" | "output" | "cancel") {
    const source: CanvasNodeData = { id: "source", type: CanvasNodeType.Text, title: "fixture", position: { x: 0, y: 0 }, width: 320, height: 240, metadata: {} };
    let nodes = [source];
    const requests: Array<{ count: string; nodeId: string }> = [];
    const consumed: string[] = [];
    const finished: string[] = [];
    const errors: string[] = [];
    const controller = new AbortController();
    const input: CanvasGenerationExecution = {
        nodeId: source.id, sourceNode: source, canvasNodes: nodes, canvasConnections: [], prompt: "fixture", effectivePrompt: "fixture", generationConfig: { ...config, ...overrides },
        generationContext: { prompt: "fixture", referenceImages: [], referenceVideos: [], referenceAudios: [], characterReferences: [], resolvedCharacterVersions: [], resolvedCharacterVoices: [], textCount: 1, imageCount: 0, videoCount: 0, audioCount: 0 },
        controller, projectId: "", editingTextNode: false, styleMetadata: {}, skillMetadata: { skillIds: [], skillVersions: [], skillFiles: [] },
        setNodes: (next) => { nodes = typeof next === "function" ? next(nodes) : next; },
        setConnections: () => {}, setSelectedNodeIds: () => {}, setSelectedConnectionId: () => {}, setDialogNodeId: () => {}, registerPendingNodeIds: () => {},
        startGenerationRequest: () => controller, finishGenerationRequest: (id) => { finished.push(id); }, showError: (error) => { errors.push(error); },
        bindGenerationTask: (id, updated) => { nodes = nodes.map((node) => node.id === id ? { ...node, metadata: { ...node.metadata, taskId: updated.id } } : node); },
        applyGenerationTaskResult: async (id, updated) => {
            const node = nodes.find((item) => item.id === id)!;
            if (failure === "output" && node.metadata?.taskOutputIndex === 1) throw new Error("output unavailable");
            nodes = (await applyGenerationTaskResultToNodes(nodes, updated, id)).nodes;
            consumed.push(id);
        },
    };
    await executeImageGeneration(input, async (request, dependencies) => {
        requests.push({ count: request.config.count, nodeId: request.nodeId });
        if (failure === "submit") throw new Error("upstream failed");
        if (failure === "cancel") { controller.abort(); throw new DOMException("Aborted", "AbortError"); }
        dependencies.bindTask(task);
        await dependencies.consumeTask(task);
        return JSON.parse(task.resultJson!);
    });
    return { nodes, requests, consumed, finished, errors };
}

describe("canvas native image batch", () => {
    test("one request with n=3 attaches distinct indexed results to all children", async () => {
        const result = await execute();
        expect(result.requests).toHaveLength(1);
        expect(result.requests[0]!.count).toBe("3");
        const children = result.nodes.filter((node) => node.metadata?.batchRootId);
        expect(children.map((node) => node.metadata?.storageKey)).toEqual(["resource:image-0", "resource:image-1", "resource:image-2"]);
        expect(children.map((node) => node.metadata?.taskOutputIndex)).toEqual([0, 1, 2]);
        expect(children.every((node) => node.metadata?.taskSubmissionNodeId === result.requests[0]!.nodeId && node.metadata.taskId === task.id)).toBe(true);
        expect(result.consumed).toHaveLength(3);
        expect(result.finished).toHaveLength(4);
        expect(result.errors).toEqual([]);
    });
    test("default mode still submits three independent count=1 tasks", async () => {
        const result = await execute({ generationOptions: {} });
        expect(result.requests.map((request) => request.count)).toEqual(["1", "1", "1"]);
        expect(result.nodes.filter((node) => node.metadata?.batchRootId).every((node) => node.metadata?.taskOutputIndex === undefined)).toBe(true);
    });
    test("one image stays one task with output zero", async () => {
        const result = await execute({ count: "1" });
        expect(result.requests.map((request) => request.count)).toEqual(["1"]);
        expect(result.consumed).toHaveLength(1);
    });
    test("one output failure does not hide successful siblings or resubmit", async () => {
        const result = await execute({}, "output");
        expect(result.requests).toHaveLength(1);
        expect(result.consumed).toHaveLength(2);
        expect(result.nodes.find((node) => node.metadata?.isBatchRoot)?.metadata).toMatchObject({ status: "success", batchFailedCount: 1 });
    });
    test("submission failure fails all children without automatic retry", async () => {
        const result = await execute({}, "submit");
        expect(result.requests).toHaveLength(1);
        expect(result.nodes.filter((node) => node.metadata?.batchRootId).every((node) => node.metadata?.status === "error")).toBe(true);
    });
    test("cancellation finishes every request handle without resubmitting", async () => {
        const result = await execute({}, "cancel");
        expect(result.requests).toHaveLength(1);
        expect(result.finished).toHaveLength(4);
        expect(result.errors).toEqual([]);
    });
    test("logical contract is supported while stale options and workflows are excluded", () => {
        expect(nativeImageBatchEnabled(config)).toBe(true);
        expect(nativeImageBatchEnabled({ ...config, taskWorkflowProvider: "comfyui" })).toBe(false);
        expect(nativeImageBatchEnabled({ ...config, model: "other" })).toBe(false);
        expect(nativeImageBatchEnabled({ ...config, generationOptions: { nativeBatch: true, contractKind: "wan2.7-i2v" } })).toBe(false);
        expect(nativeImageBatchEnabled({ ...config, model: "system::opaque", channels: [{ ...defaultConfig.channels[0]!, id: "system", models: ["opaque"], modelCosts: [{ model: "opaque", billingMode: "fixed_request", unitPriceMicrocredits: 1, capabilityConfig: { generationContract: "qwen-image-3.0-pro" } }] }] })).toBe(true);
    });
    test("output mapping survives JSON restore and never overwrites a sibling for a deleted target", async () => {
        const { nodes } = await execute();
        const restored: CanvasNodeData[] = JSON.parse(JSON.stringify(nodes));
        const child = restored.find((node) => node.metadata?.taskOutputIndex === 2)!;
        const result = await applyGenerationTaskResultToNodes(restored, task, child.id);
        expect(result.node?.metadata?.storageKey).toBe("resource:image-2");
        expect((await applyGenerationTaskResultToNodes(restored.filter((node) => node.id !== child.id), task, child.id)).updated).toBe(false);
    });
    test("missing indexed output rejects instead of silently attaching the first image", async () => {
        const { nodes } = await execute();
        const child = nodes.find((node) => node.metadata?.taskOutputIndex === 2)!;
        await expect(applyGenerationTaskResultToNodes(nodes, { ...task, resultJson: JSON.stringify({ mode: "image", images: [{ dataUrl: "first" }] }) }, child.id)).rejects.toThrow("没有返回图片");
    });
    test("per-output retries have distinct idempotency keys, stable for repeated clicks", async () => {
        const first = await createGenerationRetryContext(task.id, `${task.id}:output:0`);
        const second = await createGenerationRetryContext(task.id, `${task.id}:output:1`);
        expect(first.clientOperationId).not.toBe(second.clientOperationId);
        expect(first).toEqual(await createGenerationRetryContext(task.id, `${task.id}:output:0`));
    });
    test("node-local batch options reach the quote configuration after restoring a draft", () => {
        const node: CanvasNodeData = { id: "node", type: CanvasNodeType.Image, title: "", position: { x: 0, y: 0 }, width: 320, height: 240, metadata: { model: config.model, count: 3, generationOptions: config.generationOptions } };
        const channel = createModelChannel({ id: "qwen", models: [config.model], interfaceType: "dashscope-image" });
        const resolved = buildNodeConfig({ ...config, channels: [channel] }, node, "image", { capability: "image" });
        expect(resolved.generationOptions).toEqual(config.generationOptions);
        expect(nativeImageBatchEnabled(resolved)).toBe(true);
    });
    test("materialized outputs attach their own asset and acknowledge each output idempotently", async () => {
        const { nodes } = await execute();
        const child = nodes.find((node) => node.metadata?.taskOutputIndex === 1)!;
        const asset: ImageAsset = { id: "materialized-1", kind: "image", title: "fixture", coverUrl: "", tags: [], createdAt: "", updatedAt: "", data: { dataUrl: "fixture", storageKey: "resource:materialized-1", width: 512, height: 512, bytes: 1, mimeType: "image/png" } };
        const originalAssets = useAssetStore.getState().assets;
        useAssetStore.setState({ assets: [asset] });
        try {
            const output = { outputIndex: 1, mediaType: "image" as const, materializedAssetId: asset.id };
            const applied = await applyMaterializedGenerationTaskResultToNodes(nodes, task, output, "effect-1", child.id);
            expect(applied.node?.metadata).toMatchObject({ storageKey: "resource:materialized-1", assetId: asset.id, generationEffectKeys: ["effect-1"] });
            const repeated = await applyMaterializedGenerationTaskResultToNodes(applied.nodes, task, output, "effect-1", child.id);
            expect(repeated.nodes).toBe(applied.nodes);
            await expect(applyMaterializedGenerationTaskResultToNodes(nodes, task, { ...output, outputIndex: 0 }, "wrong-slot", child.id)).rejects.toThrow("序号");
        } finally { useAssetStore.setState({ assets: originalAssets }); }
    });
});
