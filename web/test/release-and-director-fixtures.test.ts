import { expect, test } from "bun:test";
import { readFileSync } from "node:fs";
import { directorE2EFixture } from "../scripts/director-e2e-fixtures.mjs";

test("Director fixtures cover exact read-only dependencies without hiding other failures", () => {
    expect(JSON.parse(directorE2EFixture("/api/public/appearance")!.body).code).toBe(0);
    expect(directorE2EFixture("/favicon.ico")!.contentType).toBe("image/svg+xml");
    expect(directorE2EFixture("/api/tasks")).toBeUndefined();
    expect(directorE2EFixture("/api/public/appearance", "POST")).toBeUndefined();
    const config = readFileSync(new URL("../vite.config.ts", import.meta.url), "utf8");
    expect(config).toContain('process.env.CANVAS_DIRECTOR_E2E === "1"');
    expect(config).toContain('fileName: "build-info.json"');
    expect(config).toContain('"import.meta.env.VITE_BUILD_COMMIT"');
});

test("fork display version maps to a Docker-compatible distinct tag", () => {
    const version = readFileSync(new URL("../../VERSION", import.meta.url), "utf8").trim();
    const tag = version.replace(/^v/, "").replaceAll("+", "-");
    expect(version).toBe("v1.2.7+dingxin.1");
    expect(tag).toBe("1.2.7-dingxin.1");
    expect(tag).toMatch(/^[a-zA-Z0-9_][a-zA-Z0-9_.-]{0,127}$/);
});
