import { readdirSync } from "node:fs";
import { spawnSync } from "node:child_process";

const files = ["test", "src"]
    .flatMap((directory) => readdirSync(new URL(`../${directory}/`, import.meta.url), { recursive: true }).map((name) => `${directory}/${name.replaceAll("\\", "/")}`))
    .filter((name) => /\.(test|spec)\.(ts|tsx|js|jsx|mjs|cjs)$/.test(name))
    .sort();
// Other suites register built-in slots at module initialization. Keep the
// registry's empty-initial-state contract in a separate process, not disabled.
const isolated = ["test/editor-slot-registry.test.ts"];
for (const group of [files.filter((name) => !isolated.includes(name)), ...isolated.map((name) => [name])]) {
    const result = spawnSync(process.execPath, ["scripts/run-with-panic-guard.mjs", "--", process.execPath, "test", ...group.map((name) => `./${name}`)], { stdio: "inherit" });
    if (result.error || result.status !== 0) process.exit(result.status || 1);
}
