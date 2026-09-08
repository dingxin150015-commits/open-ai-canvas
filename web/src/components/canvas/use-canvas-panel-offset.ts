import { useLayoutEffect, useState, type RefObject } from "react";
import { canvasPanelLeft } from "@/lib/canvas/canvas-chrome-layout";

export function useCanvasPanelOffset(open: boolean, rootRef: RefObject<HTMLElement | null>, width: number) {
    const [left, setLeft] = useState(0);
    useLayoutEffect(() => {
        const root = rootRef.current;
        if (!open || !root) return;
        const update = () => setLeft(canvasPanelLeft(root.getBoundingClientRect().left, width, window.innerWidth));
        update();
        const observer = new ResizeObserver(update);
        observer.observe(root.closest("[data-canvas-bottom-bar]") || root);
        window.addEventListener("resize", update);
        return () => {
            observer.disconnect();
            window.removeEventListener("resize", update);
        };
    }, [open, rootRef, width]);
    return left;
}
