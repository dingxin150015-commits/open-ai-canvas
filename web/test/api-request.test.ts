import { describe, expect, test } from "bun:test";

import { quoteModelCatalog } from "../src/services/api/logical-models";
import { ApiError, apiClient, request } from "../src/services/api/request";

describe("backend API request error semantics", () => {
    test("unwraps a successful backend envelope", async () => {
        await expect(request(Promise.resolve({ data: { code: 0, data: { id: "task-1" }, msg: "ok" }, status: 200 }))).resolves.toEqual({ id: "task-1" });
    });

    test("preserves business code from a non-zero envelope", async () => {
        const thrown = await request(Promise.resolve({ data: { code: 40901, data: null, msg: "任务状态已变化" }, status: 200 })).catch((error) => error);

        expect(thrown).toBeInstanceOf(ApiError);
        expect(thrown).toMatchObject({ name: "ApiError", status: 200, code: 40901, message: "任务状态已变化", retryable: false });
    });

    test("preserves HTTP status, backend code and retryability", async () => {
        const axiosError = {
            isAxiosError: true,
            message: "Request failed with status code 429",
            response: {
                status: 429,
                data: {
                    code: 42901,
                    data: null,
                    msg: "请求过于频繁，请稍后重试",
                    errorCode: "request_throttled",
                    errorCategory: "quota",
                    retryable: true,
                    requestId: "req_12345678",
                },
                headers: { "x-request-id": "header-request-id" },
            },
        };
        const thrown = await request(Promise.reject(axiosError)).catch((error) => error);

        expect(thrown).toBeInstanceOf(ApiError);
        expect(thrown).toMatchObject({ status: 429, code: 42901, errorCode: "request_throttled", errorCategory: "quota", requestId: "req_12345678", message: "请求过于频繁，请稍后重试", retryable: true });
        expect(thrown.cause).toBe(axiosError);
    });

    test("uses the response request-id header when an older backend body has no diagnostic id", async () => {
        const axiosError = {
            isAxiosError: true,
            message: "Request failed with status code 500",
            response: {
                status: 500,
                data: { code: 500, data: null, msg: "系统处理失败，请稍后重试" },
                headers: { "x-request-id": "req_header_12345678" },
            },
        };
        const thrown = await request(Promise.reject(axiosError)).catch((error) => error);

        expect(thrown).toMatchObject({ requestId: "req_header_12345678", retryable: true });
    });

    test("converts Axios cancellation to AbortError", async () => {
        const thrown = await request(Promise.reject({ __CANCEL__: true, code: "ERR_CANCELED" })).catch((error) => error);

        expect(thrown).toMatchObject({ name: "AbortError", message: "请求已取消" });
    });
});

describe("统一模型目录报价", () => {
    test("向统一端点提交模型 ID 和完整意图", async () => {
        const originalAdapter = apiClient.defaults.adapter;
        apiClient.defaults.adapter = async (config) => {
            expect(config.url).toBe("/model-catalog/quote");
            expect(config.method).toBe("post");
            expect(JSON.parse(String(config.data))).toEqual({
                modelId: "channel-model-1",
                intent: {
                    capability: "video",
                    operation: "text_to_video",
                    inputs: { image: 0 },
                    options: { vquality: "480p", videoSeconds: 2 },
                },
            });
            return {
                data: { code: 0, data: { quote: { modelId: "channel-model-1", billingMode: "per_second", quantity: 2, amountMicrocredits: 2_000_000, estimated: false } }, msg: "ok" },
                status: 200,
                statusText: "OK",
                headers: {},
                config,
            };
        };
        try {
            const result = await quoteModelCatalog("channel-model-1", {
                capability: "video",
                operation: "text_to_video",
                inputs: { image: 0 },
                options: { vquality: "480p", videoSeconds: 2 },
            });
            expect(result.quote).toMatchObject({ modelId: "channel-model-1", quantity: 2, amountMicrocredits: 2_000_000 });
        } finally {
            apiClient.defaults.adapter = originalAdapter;
        }
    });
});
