#!/usr/bin/env python3
"""Search the local Bailian official documentation mirror with file and heading context."""

from __future__ import annotations

import argparse
import re
import sys
from pathlib import Path


PROJECT_ROOT = Path(__file__).resolve().parents[4]
DEFAULT_DOCS_ROOT = PROJECT_ROOT / "百炼千问文档" / "百炼千问文档"


def read_text(path: Path) -> str:
    raw = path.read_bytes()
    for encoding in ("utf-8-sig", "utf-8", "gb18030"):
        try:
            return raw.decode(encoding)
        except UnicodeDecodeError:
            continue
    return raw.decode("utf-8", errors="replace")


def main() -> int:
    if hasattr(sys.stdout, "reconfigure"):
        sys.stdout.reconfigure(encoding="utf-8", errors="replace")
    parser = argparse.ArgumentParser()
    parser.add_argument("query", nargs="+", help="Search terms; every term must match")
    parser.add_argument("--docs-root", type=Path, default=DEFAULT_DOCS_ROOT)
    parser.add_argument("--limit", type=int, default=30)
    args = parser.parse_args()
    terms = [term.casefold() for term in args.query]
    matches = 0
    for path in sorted(args.docs_root.rglob("*")):
        if not path.is_file() or path.suffix.lower() not in {".md", ".txt"}:
            continue
        text = read_text(path)
        folded = text.casefold()
        if not all(term in folded for term in terms):
            continue
        lines = text.splitlines()
        title = next((line[2:].strip() for line in lines if line.startswith("# ")), path.stem)
        print(f"{path.relative_to(args.docs_root).as_posix()} | {title}")
        shown = 0
        for number, line in enumerate(lines, start=1):
            if any(term in line.casefold() for term in terms):
                print(f"  {number}: {line[:240]}")
                shown += 1
                if shown >= 5:
                    break
        matches += 1
        if matches >= args.limit:
            break
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
