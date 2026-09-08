import { afterEach, expect, test } from "bun:test";
import { readFileSync } from "node:fs";

import { createThemeStore } from "../src/stores/use-theme-store";

const source = (path: string) => readFileSync(new URL(`../src/${path}`, import.meta.url), "utf8");
const saved = new Map<string, unknown>();
const useThemeStore = createThemeStore({
    getItem: (name) => saved.get(name) as { state: { theme: "light" | "dark" } } | null,
    setItem: (name, value) => {
        saved.set(name, value);
    },
    removeItem: (name) => {
        saved.delete(name);
    },
});

afterEach(() => {
    useThemeStore.getState().setCanvasTheme(null);
    saved.clear();
});

test("opening a dark canvas cannot overwrite the light workspace preference", async () => {
    useThemeStore.getState().setTheme("light");
    useThemeStore.getState().setCanvasTheme("dark");
    expect(useThemeStore.getState().theme).toBe("dark");
    expect(saved.get("infinite-canvas:theme_store")).toMatchObject({ state: { theme: "light" } });
    await useThemeStore.persist.rehydrate();
    expect(useThemeStore.getState().theme).toBe("dark");
    useThemeStore.getState().setCanvasTheme(null);
    expect(useThemeStore.getState().theme).toBe("light");
});

test("a user preference changed during a canvas session survives leaving it", () => {
    useThemeStore.getState().setTheme("dark");
    useThemeStore.getState().setCanvasTheme("dark");
    useThemeStore.getState().setTheme("light");
    useThemeStore.getState().setCanvasTheme(null);
    expect(useThemeStore.getState().theme).toBe("light");
});

test("invalid themes cannot corrupt the workspace or canvas display", () => {
    useThemeStore.getState().setTheme("light");
    useThemeStore.getState().setCanvasTheme("invalid" as "dark");
    useThemeStore.getState().setTheme("invalid" as "dark");
    expect(useThemeStore.getState().theme).toBe("light");
});

test("canvas appearance reaches the renderer and grid strength is not multiplied twice", () => {
    const page = source("pages/canvas/project.tsx");
    const renderer = page.slice(page.indexOf("<InfiniteCanvas"), page.indexOf("graphicsLayer={", page.indexOf("<InfiniteCanvas")));
    expect(renderer).toContain("appearance={canvasAppearance}");
    expect(source("components/canvas/infinite-canvas.tsx")).not.toContain('opacity: mode === "dots"');
    expect(page).toContain("setCanvasTheme(null)");
    expect(source("pages/canvas/use-canvas-project-lifecycle.ts")).not.toContain("getState().setTheme(");
});
