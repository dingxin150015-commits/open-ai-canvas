#!/usr/bin/env python3
"""Build a traceable index for the local Bailian official documentation mirror."""

from __future__ import annotations

import argparse
import hashlib
import json
import re
import sys
from collections import Counter
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[4]
DEFAULT_DOCS_ROOT = PROJECT_ROOT / "百炼千问文档" / "百炼千问文档"
DEFAULT_SOURCE_LABEL = "<project-root>/百炼千问文档/百炼千问文档"
TEXT_SUFFIXES = {".md", ".txt"}
URL_RE = re.compile(r"https?://[^\s)\]>'\"]+")
MODEL_RE = re.compile(
    r"\b(?:qwen|wan|cosyvoice|paraformer|fun-asr|sambert|gte|text-embedding|"
    r"dashscope|deepseek|glm|kimi|minimax|step|mimo|vidu)[a-zA-Z0-9._:+/-]*\b",
    re.IGNORECASE,
)


def read_text(path: Path) -> str:
    raw = path.read_bytes()
    for encoding in ("utf-8-sig", "utf-8", "gb18030"):
        try:
            return raw.decode(encoding)
        except UnicodeDecodeError:
            continue
    return raw.decode("utf-8", errors="replace")


def first_title(text: str, fallback: str) -> str:
    for line in text.splitlines():
        match = re.match(r"^#\s+(.+?)\s*$", line)
        if match:
            return match.group(1).strip()
    return fallback


def headings(text: str) -> list[str]:
    result: list[str] = []
    for line in text.splitlines():
        match = re.match(r"^#{1,3}\s+(.+?)\s*$", line)
        if match:
            result.append(match.group(1).strip())
    return result[:24]


def mode_flags(text: str) -> list[str]:
    lower = text.lower()
    flags: list[str] = []
    checks = (
        ("sync", ("同步调用", "synchronous")),
        ("async", ("异步调用", "async", "task_id", "task-id")),
        ("stream", ("stream=true", "流式输出", "websocket")),
        ("openai-compatible", ("openai兼容", "openai-compatible", "/compatible-mode/")),
        ("dashscope-native", ("dashscope", "百炼api", "百炼 api")),
    )
    for name, needles in checks:
        if any(needle in lower for needle in needles):
            flags.append(name)
    return flags


def main() -> int:
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8", errors="replace")
    parser = argparse.ArgumentParser()
    parser.add_argument("--docs-root", type=Path, default=DEFAULT_DOCS_ROOT)
    parser.add_argument("--output-dir", type=Path, required=True)
    args = parser.parse_args()

    docs_root = args.docs_root.resolve()
    output_dir = args.output_dir.resolve()
    source_label = DEFAULT_SOURCE_LABEL if docs_root == DEFAULT_DOCS_ROOT.resolve() else str(docs_root)
    if not docs_root.is_dir():
        raise SystemExit(f"Documentation root not found: {docs_root}")
    output_dir.mkdir(parents=True, exist_ok=True)

    entries: list[dict[str, object]] = []
    top_counts: Counter[str] = Counter()
    group_counts: Counter[str] = Counter()
    total_lines = 0

    files = sorted(
        path for path in docs_root.rglob("*")
        if path.is_file() and path.suffix.lower() in TEXT_SUFFIXES
    )
    for path in files:
        raw = path.read_bytes()
        text = read_text(path)
        relative = path.relative_to(docs_root).as_posix()
        parts = relative.split("/")
        top = parts[0] if len(parts) > 1 else "root"
        group = "/".join(parts[:2]) if len(parts) > 2 else top
        line_count = len(text.splitlines())
        total_lines += line_count
        top_counts[top] += 1
        group_counts[group] += 1
        urls = sorted(set(URL_RE.findall(text)))
        models = sorted(
            set(match.group(0) for match in MODEL_RE.finditer(text)),
            key=lambda value: (value.casefold(), value),
        )
        entries.append(
            {
                "path": relative,
                "title": first_title(text, path.stem),
                "sha256": hashlib.sha256(raw).hexdigest(),
                "bytes": len(raw),
                "lines": line_count,
                "headings": headings(text),
                "models": models[:80],
                "urls": urls[:80],
                "flags": mode_flags(text),
            }
        )

    catalog = {
        "source_root": source_label,
        "document_count": len(entries),
        "line_count": total_lines,
        "documents": entries,
    }
    with (output_dir / "source-catalog.json").open("w", encoding="utf-8", newline="\n") as catalog_file:
        catalog_file.write(json.dumps(catalog, ensure_ascii=False, indent=2) + "\n")

    lines = [
        "# 百炼官方文档逐篇索引",
        "",
        f"- 源目录：`{source_label}`",
        f"- 文档数：{len(entries)}",
        f"- 总行数：{total_lines}",
        "- 每个文件均已完整读取并记录 SHA-256；摘要不替代原文。",
        "",
        "## 顶层分类",
        "",
        "| 分类 | 文件数 |",
        "| --- | ---: |",
    ]
    for name, count in sorted(top_counts.items()):
        lines.append(f"| `{name}` | {count} |")
    lines.extend(["", "## 二级分类", "", "| 分类 | 文件数 |", "| --- | ---: |"])
    for name, count in sorted(group_counts.items()):
        lines.append(f"| `{name}` | {count} |")
    lines.extend([
        "",
        "## 全部文档",
        "",
        "| 路径 | 标题 | 行数 | 模式 | SHA-256 |",
        "| --- | --- | ---: | --- | --- |",
    ])
    for entry in entries:
        title = str(entry["title"]).replace("|", "\\|").replace("\n", " ")
        flags = ", ".join(entry["flags"]) or "-"
        lines.append(
            f"| `{entry['path']}` | {title} | {entry['lines']} | {flags} | `{entry['sha256']}` |"
        )
    with (output_dir / "source-index.md").open("w", encoding="utf-8", newline="\n") as index_file:
        index_file.write("\n".join(lines) + "\n")
    print(json.dumps({"documents": len(entries), "lines": total_lines, "output": str(output_dir)}, ensure_ascii=False))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
