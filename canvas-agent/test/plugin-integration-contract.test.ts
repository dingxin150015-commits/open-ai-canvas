import assert from "node:assert/strict";
import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";
import { test } from "node:test";

import { CANVAS_GENERATION_CONTINUATION_TIMEOUT_MS } from "../src/canvas-tool-timeouts.js";

const repositoryRoot = path.resolve(path.dirname(fileURLToPath(import.meta.url)), "../..");
const pluginRoot = path.join(repositoryRoot, "plugins", "yingce");

function read(relativePath: string) {
    return fs.readFileSync(path.join(repositoryRoot, relativePath), "utf8");
}

test("external plugin timeout covers the Canvas generation continuation window", () => {
    const config = JSON.parse(read("plugins/yingce/.mcp.json")) as {
        mcpServers: { yingce: { tool_timeout_sec: number } };
    };
    const timeoutMs = config.mcpServers.yingce.tool_timeout_sec * 1_000;
    assert.ok(timeoutMs > CANVAS_GENERATION_CONTINUATION_TIMEOUT_MS);
});

test("open-canvas docs use Vite and the signed Local Runtime handshake without URL tokens", () => {
    for (const relativePath of ["plugins/yingce/README.md", "plugins/yingce/skills/open-canvas/SKILL.md"]) {
        const content = read(relativePath);
        assert.doesNotMatch(content, /[?&](?:agentToken|agentUrl)=|CANVAS_URL\s*=|bunx\s+next|\bnext dev\b/i, relativePath);
        assert.match(content, /FRAMEFIELD_TRUSTED_WEB_ORIGINS/, relativePath);
        assert.match(content, /\/canvas\?mode=new/, relativePath);
    }
    assert.match(read("plugins/yingce/skills/open-canvas/SKILL.md"), /bun run dev -- --host 127\.0\.0\.1 --port/);
});

test("all runtime Backend Compose entrypoints pass the provider plugin switch", () => {
    for (const filename of [
        "docker-compose.yml",
        "docker-compose.local.yml",
        "docker-compose.dev.yml",
        "docker-compose.deploy.yml",
        "docker-compose.server.yml",
    ]) {
        assert.match(read(filename), /ENABLE_PROVIDER_PLUGINS:\s+\$\{ENABLE_PROVIDER_PLUGINS:-false\}/, filename);
    }
});

test("plugin manifest and referenced files remain complete", () => {
    const manifest = JSON.parse(read("plugins/yingce/.codex-plugin/plugin.json")) as {
        name: string;
        version: string;
        skills: string;
        mcpServers: string;
    };
    assert.equal(manifest.name, path.basename(pluginRoot));
    assert.match(manifest.version, /^\d+\.\d+\.\d+$/);
    assert.equal(fs.existsSync(path.join(pluginRoot, manifest.skills)), true);
    assert.equal(fs.existsSync(path.join(pluginRoot, manifest.mcpServers)), true);
});
