import { expect, test } from "bun:test";

test("Docker web image builds every Comfy Bridge download instead of relying on host artifacts", async () => {
    const dockerfile = await Bun.file(new URL("../../Dockerfile", import.meta.url)).text();
    const ignored = await Bun.file(new URL("../../.dockerignore", import.meta.url)).text();
    expect(dockerfile).toContain("AS comfy-bridge-build");
    expect(dockerfile).toContain("COPY canvas-agent/native/comfy-bridge ./");
    for (const target of ["OpenAICanvas-ComfyBridge.exe", "OpenAICanvas-ComfyBridge-linux-amd64", "OpenAICanvas-ComfyBridge-linux-arm64"]) {
        expect(dockerfile).toContain(`/out/${target}`);
    }
    expect(dockerfile).toContain("COPY --from=comfy-bridge-build /out/ /app/web/public/");
    expect(ignored).toContain("/web/public/OpenAICanvas-ComfyBridge-linux-amd64");
    expect(ignored).toContain("/web/public/OpenAICanvas-ComfyBridge-linux-arm64");
});
