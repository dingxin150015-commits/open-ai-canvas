import type { ChannelModel, ChannelModelFetchResult, ChannelModelSupportStatus } from "@/services/api/wallet";

export function channelModelSupportMeta(status: ChannelModelSupportStatus) {
    return (
        {
            unknown: { label: "待配置", tone: "info" as const },
            ready: { label: "可用", tone: "success" as const },
            planned: { label: "计划支持", tone: "info" as const },
            unsupported: { label: "暂不支持", tone: "warning" as const },
            deprecated: { label: "已废弃", tone: "neutral" as const },
        }[status] || { label: "待配置", tone: "info" as const }
    );
}

export function channelModelExecutionMeta(item: ChannelModel) {
    const state = item.eligibility;
    if (!state) return channelModelSupportMeta(item.supportStatus);
    if (state.executionSupported) return { label: "协议就绪", tone: "success" as const };
    if (!state.configurationValid) return { label: "待配置", tone: "info" as const };
    return { label: "执行受限", tone: "warning" as const };
}

export function channelModelExecutionStatus(item: ChannelModel): ChannelModelSupportStatus {
    if (!item.eligibility) return item.supportStatus;
    if (item.eligibility.executionSupported) return "ready";
    if (!item.eligibility.configurationValid) return "unknown";
    return item.supportStatus === "unsupported" || item.supportStatus === "deprecated" ? item.supportStatus : "planned";
}

export function isChannelModelReadOnly(item: ChannelModel | null | undefined) {
    return Boolean(item && item.eligibility?.canEditConfiguration === false);
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
