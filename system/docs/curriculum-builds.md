# Curriculum builds and maintenance

Editable curricula live under `programs/<id>/curriculum/`, beside their viewing PDFs.
The shared style is [curriculum.sty](../latex/curriculum.sty).
Books, Markdown sources, and external URLs do not need PDF compilation.

Use TeX Live 2026 with pdfTeX 3.141592653-2.6-1.40.29, latexmk 4.88, sourcesans 3.0, and Latin Modern fonts.
Required packages are geometry, fontenc, microtype, array, enumitem, titlesec, titletoc, fancyhdr, ragged2e, multicol, url, and hyperref.
From the repository root, build all configured curricula or one PDF:

```sh
make pdfs
make programs/algorithms/curriculum/Algorithms_Mastery_Course.pdf
```

The build replaces a viewing PDF only after successful compilation.
Render and inspect changed PDFs before delivery.
`make clean-pdfs` removes temporary `build/pdf/` files while retaining viewing PDFs and `build/bin/course`.
Learner PDFs under `sessions/` are outside this build.

## Add or maintain a program

Read the [configuration reference](course-configuration.md).
Add `programs/<id>/course.json` and `AGENTS.md`; no central registry needs updating.
Source paths are relative to the program; shared guidance paths are relative to the repository root.
Keep learner exercises and their build systems inside the program.

From the repository root, run `make build`, `make vet`, and `bash course validate` after infrastructure changes.
Exercise changed session commands with `--state /absolute/temporary/state.json` to protect progress.
Preserve curricula, viewing PDFs, learner work, and records.
Keep infrastructure tests absent and preserve learner-written exercise tests.
