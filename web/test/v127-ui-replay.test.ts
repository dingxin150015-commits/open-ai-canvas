import { expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { canvasChromeLayout, canvasPanelLeft, partitionDockEntries } from "../src/lib/canvas/canvas-chrome-layout";
import { createAutomaticRuntimeDiscovery, createLocalRuntimeStore } from "../src/stores/use-local-runtime-store";
import { canvasRuntimeRetryDelay } from "../src/lib/canvas/local-runtime-connection";
import { modelQuoteRequest, modelQuoteValidationError } from "../src/lib/model-pricing";
import { defaultConfig, createModelChannel, type AiConfig } from "../src/stores/use-config-store";
import type { CapabilitySpec } from "../src/services/api/logical-models";
import type { ModelRequirements } from "../src/lib/model-selection";

const source = (p: string) => readFileSync(new URL(`../src/${p}`, import.meta.url), "utf8");

test("actual canvas width controls compact layout and overflow never drops commands", () => {
    const commands = Array.from({ length: 18 }, (_, i) => ({ id: `tool-${i}` }));
    for (const width of [208, 300, 420, 640, 900, 1280, 1920]) {
        const layout = canvasChromeLayout(width);
        const split = partitionDockEntries(commands, layout.mainLimit);
        expect(split.visible.length).toBeLessThanOrEqual(layout.mainLimit);
        expect(new Set([...split.visible, ...split.overflow].map((x) => x.id)).size).toBe(commands.length);
    }
    expect(canvasChromeLayout(700).compactTop).toBe(true);
    expect(partitionDockEntries(commands, 1, ["tool-9"]).visible[0].id).toBe("tool-9");
});

test("popover horizontal bounds include sidebars and DevTools-sized viewports", () => {
    for (const viewport of [320, 600, 1280])
        for (const anchor of [20, viewport / 2, viewport - 35]) {
            const offset = canvasPanelLeft(anchor, 312, viewport);
            expect(anchor + offset).toBeGreaterThanOrEqual(12);
            expect(anchor + offset + Math.min(312, viewport - 24)).toBeLessThanOrEqual(viewport - 12);
        }
});

test("one non-stretching bottom row owns all controls, with explicit focus and labeled assets", () => {
    const page = source("pages/canvas/project.tsx");
    expect(page.match(/data-canvas-bottom-bar/g)).toHaveLength(1);
    for (const component of ["CanvasZoomControls", "CanvasAssetTray", "CanvasToolbar", "CanvasWorkspaceModeSwitch"]) expect(page.match(new RegExp(`<${component}\\s`, "g"))).toHaveLength(1);
    expect(page).toContain("canvasChromeLayout(size.width)");
    expect(source("hooks/use-focus-mode.ts")).not.toContain("smallScreen ||");
    expect(source("components/canvas/canvas-asset-tray.tsx")).not.toContain("-right-1.5 -top-1.5");
    expect(source("components/canvas/canvas-asset-tray.tsx")).toContain("displayLabel:");
});

test("shared auto discovery survives remounts without repeated offline probes", async () => {
    let run: (() => void) | undefined;
    let calls = 0;
    const acquire = createAutomaticRuntimeDiscovery(
        async () => {
            calls++;
        },
        (fn) => {
            run = fn;
            return () => {
                run = undefined;
            };
        },
    );
    const first = acquire();
    const second = acquire();
    run?.();
    first();
    second();
    const remount = acquire();
    run?.();
    await Promise.resolve();
    expect(calls).toBe(1);
    remount();
});

test("StrictMode cleanup before discovery cancels only its scheduled attempt", () => {
    let run: (() => void) | undefined;
    let calls = 0;
    const acquire = createAutomaticRuntimeDiscovery(
        async () => {
            calls++;
        },
        (fn) => {
            run = fn;
            return () => {
                run = undefined;
            };
        },
    );
    acquire()();
    expect(run).toBeUndefined();
    const cleanup = acquire();
    run?.();
    expect(calls).toBe(1);
    cleanup();
});

test("a concurrent connecting caller does not abort the owner handshake", async () => {
    let release!: () => void;
    let calls = 0;
    let ownerSignal: AbortSignal | undefined;
    const gate = new Promise<void>((resolve) => {
        release = resolve;
    });
    const store = createLocalRuntimeStore({
        client: {
            async connect(signal) {
                calls++;
                ownerSignal = signal;
                await gate;
                return { state: "connected", runtimeVersion: 2, session: { sessionId: "fixture", keyId: "fixture", scopes: ["runtime:status"], expiresAt: "2099-01-01T00:00:00Z" } };
            },
            async request() {
                return new Response(JSON.stringify({ ok: true, runtime: { id: "framefield-local-runtime", version: "fixture", apiVersion: 2 }, modules: [] }));
            },
        },
    });
    const first = store.getState().connect();
    const other = new AbortController();
    const second = store.getState().ensureConnected(other.signal);
    other.abort();
    await second;
    expect(ownerSignal?.aborted).toBe(false);
    release();
    await first;
    expect(calls).toBe(1);
    expect(store.getState().connection).toBe("connected");
});

test("automatic stream retries have a finite budget", () => {
    expect([1, 2, 3, 4, 5].map(canvasRuntimeRetryDelay)).toEqual([2000, 5000, 15000, null, null]);
    expect(source("components/canvas/canvas-local-agent-panel.tsx")).toContain("enabled: false");
});

function managedConfig() {
    const spec: CapabilitySpec = {
        version: 1,
        capability: "video",
        operations: ["text_to_video", "image_to_video"],
        inputs: { image: { min: 0, max: 1 }, video: { min: 0, max: 0 }, audio: { min: 0, max: 1 } },
        options: {
            size: { values: ["16:9"] },
            videoSeconds: { values: [5, 10] },
            vquality: { values: ["720p"] },
            videoGenerateAudio: { values: [true, false] },
            videoWatermark: { values: [true, false] },
        },
    };
    const channel = createModelChannel({
        id: "fixture",
        name: "fixture",
        scope: "system",
        models: ["video"],
        modelCosts: [{ model: "video", capability: "video", billingMode: "fixed_request", unitPriceMicrocredits: 0, logicalModelId: "LOGICAL_fixture", logicalCapabilitySpec: spec }],
    });
    const config: AiConfig = { ...defaultConfig, channels: [channel], model: "fixture::video", videoModel: "fixture::video", size: "16:9", vquality: "720P", videoSeconds: "5", videoGenerateAudio: "false", videoWatermark: "false" };
    const requirements: ModelRequirements = { capability: "video", input: { textCount: 1, imageCount: 0, characterCount: 0, videoCount: 0, audioCount: 0 } };
    return { config, requirements, spec };
}

test("logical duration and exact outgoing options gate quotes without rounding the shot", () => {
    const { config, requirements } = managedConfig();
    expect(modelQuoteRequest(config, config.model, "video", requirements)?.modelID).toBe("LOGICAL_fixture");
    for (const seconds of ["0", "3", "5.5", "15"]) {
        expect(modelQuoteValidationError(config, config.model, "video", { ...requirements, videoSeconds: seconds })).toContain("时长");
        expect(modelQuoteRequest(config, config.model, "video", { ...requirements, videoSeconds: seconds })).toBeUndefined();
    }
    expect(modelQuoteRequest(config, config.model, "video", { ...requirements, options: { videoSeconds: 3 } })).toBeUndefined();
    expect(config.videoSeconds).toBe("5");
});

test("logical resolution, references, wildcard and legacy channel ID retain their contracts", () => {
    const { config, requirements, spec } = managedConfig();
    expect(modelQuoteRequest({ ...config, vquality: "1080p" }, config.model, "video", requirements)).toBeUndefined();
    expect(modelQuoteRequest(config, config.model, "video", { ...requirements, input: { ...requirements.input!, imageCount: 2 } })).toBeUndefined();
    spec.options!.vquality.values = ["*"];
    expect(modelQuoteValidationError({ ...config, vquality: "1080p" }, config.model, "video", requirements)).toBe("");
    const cost = config.channels[0].modelCosts![0];
    delete cost.logicalModelId;
    delete cost.logicalCapabilitySpec;
    cost.channelModelId = "MODEL_fixture";
    expect(modelQuoteRequest(config, config.model, "video", requirements)?.modelID).toBe("MODEL_fixture");
});
