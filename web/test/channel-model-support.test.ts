import { describe, expect, test } from "bun:test";

import { channelModelFetchSummary, channelModelSupportMeta, channelModelExecutionMeta, channelModelExecutionStatus, isChannelModelReadOnly } from "../src/pages/admin/components/channel-model-support";

describe("channel model support status", () => {
    test("discovery status does not lock configuration; backend permission does", () => {
        expect(channelModelSupportMeta("ready").label).toBe("可用");
        for (const status of ["planned", "unsupported", "deprecated"] as const) {
            expect(isChannelModelReadOnly({ supportStatus: status, eligibility: { canEditConfiguration: true } } as never)).toBe(false);
        }
        expect(isChannelModelReadOnly({ supportStatus: "ready", eligibility: { canEditConfiguration: false } } as never)).toBe(true);
        expect(channelModelExecutionMeta({ supportStatus: "planned", eligibility: { executionSupported: true } } as never).label).toBe("协议就绪");
        expect(channelModelExecutionMeta({ supportStatus: "planned", eligibility: { configurationValid: false } } as never).label).toBe("待配置");
        expect(channelModelExecutionStatus({ supportStatus: "planned", eligibility: { configurationValid: false } } as never)).toBe("unknown");
        expect(channelModelExecutionStatus({ supportStatus: "unknown", eligibility: { executionSupported: true } } as never)).toBe("ready");
    });

    test("summarizes upstream, official, capability, and sync counts", () => {
        const summary = channelModelFetchSummary({
            models: [],
            added: 80,
            updated: 241,
            upstreamCount: 240,
            supplementalCount: 80,
            capabilityCounts: { text: 240, image: 2, video: 78 },
            supportStatusCounts: { ready: 0, planned: 320, unsupported: 0, deprecated: 0, unknown: 0 },
            skippedRetired: 0,
            skippedConfigured: 0,
            unchanged: 0,
            officialCatalogReady: true,
        });
        expect(summary).toContain("上游 240");
        expect(summary).toContain("官方补充 80");
        expect(summary).toContain("视频 78");
        expect(summary).toContain("补齐 241");
    });
});
