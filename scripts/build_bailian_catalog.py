#!/usr/bin/env python3
"""Generate the reviewed Bailian native model manifest from the local official docs mirror."""

from __future__ import annotations

import argparse
import json
import re
import sys
from collections import defaultdict
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[1]
DEFAULT_DOCS_ROOT = PROJECT_ROOT / "百炼千问文档" / "百炼千问文档"
DEFAULT_OUTPUT = PROJECT_ROOT / "backend" / "internal" / "provider" / "bailian" / "catalog.generated.json"
CATALOG_VERSION = "2026-08-25"
MODEL_JSON = re.compile(r'"model"\s*:\s*"([A-Za-z0-9._/-]+)"')
MODEL_ENUM = re.compile(r"^\s{10,}-\s+([A-Za-z0-9._/-]+)\s*$", re.MULTILINE)


OPERATION_BY_FOLDER = {
    "animate-anyone": ["image_to_video"],
    "animate-anyone-template-gen": ["image_to_video"],
    "emo-video": ["audio_to_video"],
    "emoji-video": ["image_to_video"],
    "happyhorse-image-to-video": ["image_to_video"],
    "happyhorse-reference-to-video": ["reference_to_video"],
    "happyhorse-text-to-video": ["text_to_video"],
    "happyhorse-video-editing": ["video_edit"],
    "kling-video-generation": ["text_to_video", "image_to_video", "reference_to_video", "video_edit"],
    "liveportrait-video": ["image_to_video"],
    "pixverse-image-to-video": ["image_to_video"],
    "pixverse-image-to-video-first-frame": ["image_to_video"],
    "pixverse-image-to-video-first-last": ["first_last_frame_to_video"],
    "pixverse-lipsync": ["lip_sync"],
    "pixverse-motioncontrol": ["motion_control"],
    "pixverse-reference-to-video": ["reference_to_video"],
    "pixverse-text-to-video": ["text_to_video"],
    "pixverse-upscale": ["video_upscale"],
    "video-retalk": ["lip_sync"],
    "video-style-transform": ["video_edit"],
    "vidu-image-to-video": ["image_to_video"],
    "vidu-image-to-video-first-frame": ["image_to_video"],
    "vidu-reference-to-video": ["reference_to_video"],
    "vidu-start-end-to-video": ["first_last_frame_to_video"],
    "vidu-text-to-video": ["text_to_video"],
    "wan-general-video-editing": ["video_edit"],
    "wan-image-to-animation": ["image_to_video"],
    "wan-image-to-video-first-frame": ["image_to_video"],
    "wan-image-to-video-first-last-frames": ["first_last_frame_to_video"],
    "wan-reference-to-video": ["reference_to_video"],
    "wan-s2v": ["audio_to_video"],
    "wan-text-to-video": ["text_to_video"],
    "wan-video-character-swap": ["video_edit"],
    "wan27-image-to-video": ["image_to_video"],
    "wan27-reference-to-video": ["reference_to_video"],
    "wan27-text-to-video": ["text_to_video"],
    "wan27-video-editing": ["video_edit"],
    "wan30-video": [
        "text_to_video",
        "image_to_video",
        "first_last_frame_to_video",
        "reference_to_video",
        "reference_audio_to_video",
        "file_to_video",
        "link_to_video",
    ],
}


DISPLAY_NAMES = {
    "wan3.0-video": "万相 3.0 全能视频",
    "wan3.0-video-prime": "万相 3.0 全能视频 Prime",
    "wan2.7-t2v": "万相 2.7 文生视频",
    "wan2.7-i2v": "万相 2.7 图生视频",
    "wan2.7-r2v": "万相 2.7 参考生视频",
    "happyhorse-1.1-t2v": "HappyHorse 1.1 文生视频",
    "happyhorse-1.1-i2v": "HappyHorse 1.1 图生视频",
    "happyhorse-1.1-r2v": "HappyHorse 1.1 参考生视频",
    "qwen-image-3.0-pro": "通义千问图像 3.0 Pro",
    "wan2.7-image-pro": "万相图像 2.7 Pro",
}


READY_VIDEO_OPERATIONS = {
    "wan3.0-video": ["text_to_video", "image_to_video", "first_last_frame_to_video", "reference_to_video", "audio_to_video"],
    "wan2.7-t2v": ["text_to_video", "audio_to_video"],
    "wan2.7-t2v-2026-06-12": ["text_to_video", "audio_to_video"],
    "wan2.7-t2v-2026-04-25": ["text_to_video", "audio_to_video"],
    "wan2.7-i2v": ["image_to_video", "first_last_frame_to_video", "extend", "audio_to_video"],
    "wan2.7-i2v-2026-04-25": ["image_to_video", "first_last_frame_to_video", "extend", "audio_to_video"],
    "wan2.7-r2v": ["reference_to_video", "image_to_video", "audio_to_video"],
    "wan2.7-r2v-2026-06-12": ["reference_to_video", "image_to_video", "audio_to_video"],
    "happyhorse-1.1-t2v": ["text_to_video"],
    "happyhorse-1.1-i2v": ["image_to_video"],
    "happyhorse-1.1-r2v": ["reference_to_video", "image_to_video"],
}


def models_from_document(text: str) -> set[str]:
    models = set(MODEL_JSON.findall(text))
    for match in re.finditer(r"\n\s{8}model:\s*\n(?P<block>.*?)(?=\n\s{8}[a-zA-Z_]+:|\Z)", text, re.DOTALL):
        models.update(MODEL_ENUM.findall(match.group("block")))
    return {model for model in models if model and model != "string"}


def main() -> int:
    parser = argparse.ArgumentParser()
    parser.add_argument("--docs-root", type=Path, default=DEFAULT_DOCS_ROOT)
    parser.add_argument("--output", type=Path, default=DEFAULT_OUTPUT)
    args = parser.parse_args()

    docs_root = args.docs_root.resolve()
    video_root = docs_root / "api-reference" / "video-generation"
    entries: dict[str, dict[str, object]] = {}
    operations: dict[str, set[str]] = defaultdict(set)
    documents: dict[str, set[str]] = defaultdict(set)

    for path in sorted(video_root.rglob("create-task.md")):
        folder = path.parent.name
        if folder not in OPERATION_BY_FOLDER:
            raise SystemExit(f"Missing operation mapping for {folder}")
        text = path.read_text(encoding="utf-8-sig")
        relative = path.relative_to(docs_root).as_posix()
        for model_id in models_from_document(text):
            operations[model_id].update(OPERATION_BY_FOLDER[folder])
            documents[model_id].add(relative)

    # The live official API page lists Prime with the same Wan 3.0 schema; the local mirror
    # predates that enum update, so keep it as an explicit reviewed addition.
    operations["wan3.0-video-prime"].update(OPERATION_BY_FOLDER["wan30-video"])
    documents["wan3.0-video-prime"].add("api-reference/video-generation/wan30-video/create-task.md")

    for model_id in sorted(operations):
        support_status = "planned"
        support_reason = "项目模型专属 Adapter 与合同测试尚未完成"
        supported_operations = sorted(operations[model_id])
        if model_id in READY_VIDEO_OPERATIONS:
            support_status = "ready"
            support_reason = "模型专属 Adapter 与离线合同测试已完成"
            supported_operations = READY_VIDEO_OPERATIONS[model_id]
            if model_id == "wan3.0-video":
                support_reason = "文本、首帧/首尾帧和参考图片/视频/音频已适配；file/link 输入仍待实现"
        entries[model_id] = {
            "id": model_id,
            "displayName": DISPLAY_NAMES.get(model_id, model_id),
            "providerModelKey": model_id,
            "capability": "video",
            "protocol": "dashscope-video",
            "supportedEndpointTypes": ["video"],
            "apiPath": "/api/v1/services/aigc/video-generation/video-synthesis",
            "supportStatus": support_status,
            "supportReason": support_reason,
            "supportedOperations": supported_operations,
            "documentationPaths": sorted(documents[model_id]),
            "catalogSource": "bailian-official-docs",
            "catalogVersion": CATALOG_VERSION,
        }

    image_entries = [
        {
            "id": "qwen-image-3.0-pro",
            "displayName": DISPLAY_NAMES["qwen-image-3.0-pro"],
            "providerModelKey": "qwen-image-3.0-pro",
            "capability": "image",
            "protocol": "dashscope-image",
            "supportedEndpointTypes": ["image"],
            "apiPath": "/api/v1/services/aigc/multimodal-generation/generation",
            "supportStatus": "planned",
            "supportReason": "Qwen Image 能力合同与 Create UI 回归尚未完成",
            "supportedOperations": ["text_to_image", "image_edit"],
            "documentationPaths": [
                "api-reference/image-generation/qwen-text-to-image.md",
                "api-reference/image-generation/qwen-image-editing.md",
            ],
            "catalogSource": "bailian-official-docs",
            "catalogVersion": CATALOG_VERSION,
        },
        {
            "id": "wan2.7-image-pro",
            "displayName": DISPLAY_NAMES["wan2.7-image-pro"],
            "providerModelKey": "wan2.7-image-pro",
            "capability": "image",
            "protocol": "dashscope-image",
            "supportedEndpointTypes": ["image"],
            "apiPath": "/api/v1/services/aigc/multimodal-generation/generation",
            "supportStatus": "planned",
            "supportReason": "Wan 2.7 Image 模型专属能力合同尚未完成",
            "supportedOperations": ["text_to_image", "image_edit"],
            "documentationPaths": ["api-reference/image-generation/wan27-image-gen-edit/create-task.md"],
            "catalogSource": "bailian-official-docs",
            "catalogVersion": CATALOG_VERSION,
        },
    ]

    payload = {
        "provider": "bailian",
        "version": CATALOG_VERSION,
        "generatedFrom": "local official Bailian documentation mirror plus reviewed live Wan 3.0 Prime enum",
        "models": sorted([*entries.values(), *image_entries], key=lambda item: str(item["id"])),
    }
    args.output.parent.mkdir(parents=True, exist_ok=True)
    with args.output.open("w", encoding="utf-8", newline="\n") as output_file:
        output_file.write(json.dumps(payload, ensure_ascii=False, indent=2) + "\n")
    print(json.dumps({"output": str(args.output), "models": len(payload["models"]), "version": CATALOG_VERSION}, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8", errors="replace")
    raise SystemExit(main())
