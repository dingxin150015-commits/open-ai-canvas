import { dirname, resolve } from "node:path";
import { readFileSync } from "node:fs";
import { fileURLToPath } from "node:url";
import react from "@vitejs/plugin-react";
import { defineConfig } from "vite";
import { directorE2EFixture } from "./scripts/director-e2e-fixtures.mjs";

const webDir = dirname(fileURLToPath(import.meta.url));
const appVersion = process.env.CANVAS_BUILD_VERSION?.trim() || readFileSync(resolve(webDir, "../VERSION"), "utf8").trim();
const appBuildCommit = process.env.CANVAS_BUILD_COMMIT?.trim() || process.env.VITE_BUILD_COMMIT?.trim() || "unknown";
const appChangelog = readFileSync(resolve(webDir, "../CHANGELOG.md"), "utf8");
const apiProxyTarget = process.env.VITE_API_PROXY_TARGET?.trim() || "http://127.0.0.1:8080";

export default defineConfig({
    plugins: [
        react(),
        ...(process.env.CANVAS_DIRECTOR_E2E === "1"
            ? [
                  {
                      name: "director-e2e-dependencies",
                      configureServer(server: import("vite").ViteDevServer) {
                          server.middlewares.use((req, res, next) => {
                              const fixture = directorE2EFixture(req.url, req.method);
                              if (!fixture) return next();
                              res.writeHead(200, { "Content-Type": fixture.contentType, "Cache-Control": "no-store" });
                              res.end(fixture.body);
                          });
                      },
                  },
              ]
            : []),
        {
            name: "canvas-build-identity",
            generateBundle() {
                this.emitFile({ type: "asset", fileName: "build-info.json", source: JSON.stringify({ version: appVersion, commit: appBuildCommit }) });
            },
        },
    ],
    define: {
        __APP_VERSION__: JSON.stringify(appVersion),
        __APP_CHANGELOG__: JSON.stringify(appChangelog),
        "import.meta.env.VITE_APP_VERSION": JSON.stringify(appVersion),
        "import.meta.env.VITE_BUILD_COMMIT": JSON.stringify(appBuildCommit),
    },
    server: {
        proxy: {
            "/api": {
                target: apiProxyTarget,
                changeOrigin: true,
                xfwd: true,
            },
            "/oauth/linuxdo/callback": {
                target: apiProxyTarget,
                changeOrigin: true,
                xfwd: true,
            },
        },
    },
    resolve: {
        alias: {
            "@": resolve(webDir, "src"),
        },
    },
    build: {
        rolldownOptions: {
            output: {
                strictExecutionOrder: true,
                codeSplitting: {
                    includeDependenciesRecursively: false,
                    minSize: 20 * 1024,
                    groups: [
                        {
                            name: "vendor-react",
                            test: /node_modules[\\/](?:react(?:-dom|-router|-router-dom)?|scheduler|zustand|use-sync-external-store|@tanstack[\\/](?:query-core|react-query))[\\/]/,
                            priority: 30,
                        },
                        {
                            name: "vendor-icons",
                            test: /node_modules[\\/](?:lucide-react|@ant-design[\\/]icons)[\\/]/,
                            priority: 20,
                            entriesAware: true,
                            entriesAwareMergeThreshold: 48 * 1024,
                        },
                        {
                            name: "vendor-antd",
                            test: /node_modules[\\/](?:antd|@ant-design|@rc-component|rc-[^\\/]+|dayjs)[\\/]/,
                            priority: 10,
                            entriesAware: true,
                            entriesAwareMergeThreshold: 80 * 1024,
                        },
                        {
                            name: "app-shared",
                            test: /[\\/]src[\\/]/,
                            priority: 5,
                            minShareCount: 2,
                            entriesAware: true,
                            entriesAwareMergeThreshold: 48 * 1024,
                        },
                    ],
                },
            },
        },
    },
});
