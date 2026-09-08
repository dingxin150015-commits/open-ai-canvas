import { create } from "zustand";
import { createJSONStorage, persist } from "zustand/middleware";

export type ThemeName = "light" | "dark";

type ThemeStore = {
    theme: ThemeName;
    preferredTheme: ThemeName;
    canvasTheme: ThemeName | null;
    setTheme: (theme: ThemeName) => void;
    setCanvasTheme: (theme: ThemeName | null) => void;
};

const VALID_THEMES: ThemeName[] = ["light", "dark"];

export function createThemeStore(storage = createJSONStorage<{ theme: ThemeName }>(() => localStorage)) {
    return create<ThemeStore>()(
        persist(
            (set) => ({
                theme: "dark",
                preferredTheme: "dark",
                canvasTheme: null,
                // 写入前校验：非法 theme 直接忽略，防止 syncRemoteUserData 等路径绕过 merge 写入坏值
                setTheme: (next) => set((state) => (VALID_THEMES.includes(next) ? { preferredTheme: next, theme: state.canvasTheme ?? next } : state)),
                // 画布外观只在页面存活期间覆盖显示；不得持久化为工作台偏好。
                setCanvasTheme: (next) => set((state) => (next === null || VALID_THEMES.includes(next) ? { canvasTheme: next, theme: next ?? state.preferredTheme } : state)),
            }),
            {
                name: "infinite-canvas:theme_store",
                storage,
                partialize: (state) => ({ theme: state.preferredTheme }),
                // 持久化恢复校验：旧版本/坏 session 写入的非法值回退到 dark，
                // 避免 canvasThemes[非法值] = undefined 触发 "reading 'node'" 崩溃
                merge: (persisted, current) => {
                    const stored = (persisted || {}) as Partial<ThemeStore>;
                    const theme = VALID_THEMES.includes(stored.theme as ThemeName) ? (stored.theme as ThemeName) : "dark";
                    return { ...current, preferredTheme: theme, theme: current.canvasTheme ?? theme };
                },
            },
        ),
    );
}

export const useThemeStore = createThemeStore();
