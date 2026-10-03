import assert from "node:assert/strict";
import { test } from "bun:test";

function ambiguousChannels(store) {
    const shared = {
        name: "Shared",
        baseUrl: "http://127.0.0.1:8000",
        apiKey: "same-key",
        apiFormat: "openai",
        models: ["same-model"],
    };
    return [store.createModelChannel({ ...shared, id: "local-enabled", allowLocalChannel: true }), store.createModelChannel({ ...shared, id: "local-disabled", allowLocalChannel: false })];
}

test("task-center forwards explicit task provider config without inferring a channel from Zustand", async () => {
    const taskCenter = await import("../src/services/api/task-center");
    const generationTask = await import("../src/services/api/generation-task");
    const requestModule = await import("../src/services/api/request");
    const store = await import("../src/stores/use-config-store");
    const userStore = await import("../src/stores/use-user-store");
    const previousWindow = globalThis.window;
    const previousConfig = store.useConfigStore.getState().config;
    userStore.useUserStore.setState({ features: { ...userStore.defaultFeatureAvailability, desktopLocalChannelsEnabled: true } });
    const channels = ambiguousChannels(store);
    const explicitConfig = {
        ...store.defaultConfig,
        channels,
        models: ["local-enabled::same-model", "local-disabled::same-model"],
        model: "local-disabled::same-model",
        textModel: "local-disabled::same-model",
    };
    store.useConfigStore.setState({ config: explicitConfig });
    globalThis.window = { location: { hostname: "127.0.0.1" }, dispatchEvent() {} };
    assert.equal(generationTask.backendProviderConfig(explicitConfig).allowLocalChannel, false);
    assert.equal(generationTask.backendProviderConfig({ ...explicitConfig, model: "local-enabled::same-model" }).allowLocalChannel, true);

    const bodies = [];
    const originalAdapter = requestModule.apiClient.defaults.adapter;
    requestModule.apiClient.defaults.adapter = async (request) => {
        bodies.push({ url: request.url, body: JSON.parse(String(request.data)) });
        return { config: request, data: { code: 0, data: { id: `task-${bodies.length}`, type: "canvas_text", status: "queued", prompt: "p", attempts: 1, createdAt: "now", updatedAt: "now" }, msg: "ok" }, headers: {}, status: 200, statusText: "OK" };
    };
    try {
        for (const model of ["local-disabled::same-model", "local-enabled::same-model"]) {
            const providerConfig = generationTask.backendProviderConfig({ ...explicitConfig, model });
            await taskCenter.createGenerationTask({ prompt: "p", model: "same-model", input: { mode: "text", config: providerConfig } });
        }
    } finally {
        requestModule.apiClient.defaults.adapter = originalAdapter;
        userStore.useUserStore.setState({ features: userStore.defaultFeatureAvailability });
        globalThis.window = previousWindow;
        store.useConfigStore.setState({ config: previousConfig });
    }

    assert.deepEqual(
        bodies.map((item) => item.url),
        ["/tasks", "/tasks"],
    );
    assert.equal(bodies[0]?.body?.input?.config?.allowLocalChannel, false);
    assert.equal(bodies[1]?.body?.input?.config?.allowLocalChannel, true);
});
