# syntax=docker/dockerfile:1.7

# Comfy Bridge 下载包由镜像构建生成，不依赖宿主机未跟踪的二进制。
FROM --platform=$BUILDPLATFORM golang:1.25-alpine AS comfy-bridge-build
WORKDIR /src/canvas-agent/native/comfy-bridge
COPY canvas-agent/native/comfy-bridge ./
RUN mkdir -p /out && \
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/OpenAICanvas-ComfyBridge.exe . && \
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -trimpath -ldflags="-s -w" -o /out/OpenAICanvas-ComfyBridge-linux-amd64 . && \
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -trimpath -ldflags="-s -w" -o /out/OpenAICanvas-ComfyBridge-linux-arm64 .

# 构建 Vite 前端产物。
FROM --platform=$BUILDPLATFORM oven/bun:1.4.0 AS web-build

WORKDIR /app/web
COPY web/package.json web/bun.lock ./
RUN --mount=type=cache,target=/root/.bun/install/cache bun install --frozen-lockfile --cache-dir=/root/.bun/install/cache
COPY VERSION /app/VERSION
COPY CHANGELOG.md /app/CHANGELOG.md
COPY README.md /app/README.md
COPY assets /app/assets
COPY web ./
COPY --from=comfy-bridge-build /out/ /app/web/public/
ARG VITE_TLDRAW_LICENSE_KEY
ARG BUILD_VERSION
ARG BUILD_COMMIT=unknown
ARG BUILD_TIME=unknown
ENV VITE_TLDRAW_LICENSE_KEY=${VITE_TLDRAW_LICENSE_KEY}
ENV CANVAS_BUILD_VERSION=${BUILD_VERSION}
ENV CANVAS_BUILD_COMMIT=${BUILD_COMMIT}
ENV CANVAS_BUILD_TIME=${BUILD_TIME}
# 生产镜像只构建云端工作台前端；Agent Runtime 在后端 Worker 中运行。
RUN bun --bun ./node_modules/vite/bin/vite.js build

# 运行镜像：nginx 托管静态前端，并在 Compose 中把 /api 转发到后端服务。
FROM nginx:1.27-alpine

COPY --from=web-build /app/web/dist /opt/canvas-release
COPY docker/canvas-web-entrypoint.sh /usr/local/bin/canvas-web-entrypoint
RUN chmod +x /usr/local/bin/canvas-web-entrypoint
COPY nginx.conf /etc/nginx/conf.d/default.conf

ENTRYPOINT ["/usr/local/bin/canvas-web-entrypoint"]

EXPOSE 3000
HEALTHCHECK --interval=30s --timeout=3s --start-period=10s --retries=3 CMD wget -qO- http://127.0.0.1:3000/ >/dev/null || exit 1
