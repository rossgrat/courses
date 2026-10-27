# /// script
# requires-python = "==3.14.8"
# dependencies = ["pymupdf==1.26.7"]
# ///
import argparse
import hashlib
from pathlib import Path

import pymupdf


root = Path(__file__).resolve().parent.parent
source = root / "Algorithms_Mastery_Course.pdf"
target = root / "docs/course-guide.txt"
parser = argparse.ArgumentParser()
parser.add_argument("--check", action="store_true")
args = parser.parse_args()

with pymupdf.open(source) as document:
    parts = [f"Source: {source.name}\nPages: {len(document)}\n"
             f"SHA256: {hashlib.sha256(source.read_bytes()).hexdigest()}\n"
             "Text and link destinations extracted without changing the source.\n"]
    for number, page in enumerate(document, 1):
        parts.append(f"\n=== PAGE {number} ===\n\n")
        lines = [" ".join(line.split()) for line in page.get_text(sort=True).splitlines() if line.strip()]
        parts.append("\n".join(lines) + "\n")
        parts.extend(f"LINK: {link['uri']}\n" for link in page.get_links() if "uri" in link)
    derived = "".join(parts)

if args.check:
    if not target.exists() or target.read_text() != derived:
        raise SystemExit("The derived course guide differs from the PDF. Run the sync command without --check.")
    print("The derived course guide matches the PDF.")
else:
    target.write_text(derived)
    print(f"Updated {target}")
