import type { ChannelModel, ChannelModelFetchResult, ChannelModelSupportStatus } from "@/services/api/wallet";

export function channelModelSupportMeta(status: ChannelModelSupportStatus) {
    return (
        {
            ready: { label: "可用", tone: "success" as const, readOnly: false },
            planned: { label: "计划支持", tone: "info" as const, readOnly: true },
            unsupported: { label: "暂不支持", tone: "warning" as const, readOnly: true },
            deprecated: { label: "已废弃", tone: "neutral" as const, readOnly: true },
        }[status] || { label: "计划支持", tone: "info" as const, readOnly: true }
    );
}

export function isChannelModelReadOnly(item: ChannelModel | null | undefined) {
    return Boolean(item && channelModelSupportMeta(item.supportStatus).readOnly);
}

export function channelModelFetchSummary(result: ChannelModelFetchResult) {
    const capability = [
        result.capabilityCounts.text ? `文本 ${result.capabilityCounts.text}` : "",
        result.capabilityCounts.image ? `图片 ${result.capabilityCounts.image}` : "",
        result.capabilityCounts.video ? `视频 ${result.capabilityCounts.video}` : "",
        result.capabilityCounts.audio ? `音频 ${result.capabilityCounts.audio}` : "",
    ]
        .filter(Boolean)
        .join("、");
    return `上游 ${result.upstreamCount}，官方补充 ${result.supplementalCount}${capability ? `；${capability}` : ""}；新增 ${result.added}，补齐 ${result.updated}`;
}
