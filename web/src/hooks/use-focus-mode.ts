import { useCallback, useState } from "react";

const FOCUS_MODE_KEY = "canvas-focus-mode-v2";

// Resizing or docking DevTools must not override an explicit workspace choice.
function readInitialPreference(): boolean {
    try {
        return window.localStorage.getItem(FOCUS_MODE_KEY) === "true";
    } catch {
        return false;
    }
}

export function useFocusMode() {
    const [userPreference, setUserPreference] = useState<boolean>(readInitialPreference);
    const focusMode = userPreference;

    const persist = useCallback((next: boolean) => {
        setUserPreference(next);
        try {
            window.localStorage.setItem(FOCUS_MODE_KEY, String(next));
        } catch {
            // 忽略 localStorage 不可用场景，专注模式仍可在本次会话生效。
        }
    }, []);

    const enterFocusMode = useCallback(() => persist(true), [persist]);
    const exitFocusMode = useCallback(() => persist(false), [persist]);
    const toggleFocusMode = useCallback(() => persist(!userPreference), [persist, userPreference]);

    return {
        focusMode,
        enterFocusMode,
        exitFocusMode,
        toggleFocusMode,
    };
}
