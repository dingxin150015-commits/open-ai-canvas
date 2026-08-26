import { describe, expect, test } from "bun:test";

import { channelModelFetchSummary, channelModelSupportMeta, isChannelModelReadOnly } from "../src/pages/admin/components/channel-model-support";

describe("channel model support status", () => {
    test("only ready models are editable", () => {
        expect(channelModelSupportMeta("ready").readOnly).toBe(false);
        for (const status of ["planned", "unsupported", "deprecated"] as const) {
            expect(channelModelSupportMeta(status).readOnly).toBe(true);
            expect(isChannelModelReadOnly({ supportStatus: status } as never)).toBe(true);
        }
    });

    test("summarizes upstream, official, capability, and sync counts", () => {
        const summary = channelModelFetchSummary({
            models: [],
            added: 80,
            updated: 241,
            upstreamCount: 240,
            supplementalCount: 80,
            capabilityCounts: { text: 240, image: 2, video: 78 },
            supportStatusCounts: { ready: 0, planned: 320, unsupported: 0, deprecated: 0 },
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
