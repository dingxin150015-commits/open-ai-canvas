// Breakpoints use the actual canvas width after sidebar/Agent columns, not the window.
export function canvasChromeLayout(width: number) {
    if (width >= 1100) return { compact: false, mainLimit: 8, zoomLimit: 7, compactTop: false };
    if (width >= 760) return { compact: true, mainLimit: 7, zoomLimit: 3, compactTop: true };
    if (width >= 480) return { compact: true, mainLimit: 3, zoomLimit: 1, compactTop: true };
    return { compact: true, mainLimit: width >= 300 ? 2 : 0, zoomLimit: 1, compactTop: true };
}

export function partitionDockEntries<T extends { id: string; kind?: string }>(items: T[], limit: number, preferred: string[] = []) {
    const commands = items.filter((item) => item.kind !== "separator");
    const rank = (id: string) => (preferred.includes(id) ? preferred.indexOf(id) : preferred.length);
    const ranked = [...commands].sort((a, b) => rank(a.id) - rank(b.id));
    const visibleIds = new Set(ranked.slice(0, Math.max(0, limit)).map((item) => item.id));
    return { visible: commands.filter((item) => visibleIds.has(item.id)), overflow: commands.filter((item) => !visibleIds.has(item.id)) };
}

export function canvasPanelLeft(anchorLeft: number, panelWidth: number, viewportWidth: number) {
    const width = Math.min(panelWidth, Math.max(0, viewportWidth - 24));
    return Math.max(12 - anchorLeft, Math.min(0, viewportWidth - 12 - anchorLeft - width));
}
