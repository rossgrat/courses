# Courses

This repository holds learning programs with separate sources, agent roles, and progress records.

- [Algorithms](programs/algorithms/)
- [Trading Engineering](programs/trading-engineer/)
- [Oxide Systems Engineering](programs/oxide/)
- [Music Technology Engineering](programs/music-technology/)

## Layout

- `programs/<id>/` contains each program's configuration, instructions, sources, learner work, and optional state.
- `system/cli/` contains the shared Go command and its pinned toolchain version.
- `system/guidance/` contains teaching guidance selected by course configuration.
- `system/latex/` contains the shared curriculum style.
- `system/docs/` documents the configuration format.
- `build/` contains ignored binaries and temporary PDF build files.

The root `course` file launches `build/bin/course`.
The root `AGENTS.md` defines shared ownership, session, and record rules.
Each program's `course.json` selects its learning policy.
Each program's `AGENTS.md` adds subject-specific guidance.

## Setup

Install Go 1.27.3, as pinned in `system/cli/.go-version`.
Select that version in your shell before running Make targets.
From the repository root, build the command:

```sh
make build
```

The command uses the Go standard library and includes time-zone data.
After pulling infrastructure changes, run `make build` again.

Open the selected folder under `programs/` for course study.
The shared hooks select that program from the working directory, including deeper exercise folders.
A separate nested Git repository does not inherit course selection.
Opening the repository root injects no course context.

## Find a program

From the repository root, list programs and validate their configuration:

```sh
bash course list
bash course validate
```

Both commands read `programs/*/course.json` without creating state or starting sessions.
An invalid configuration stops the command with an error.

Read a selected program's policy and progress:

```sh
bash course context --course algorithms
```

From any folder inside a program, omit the identifier:

```sh
bash "$(git rev-parse --show-toplevel)/course" context
```

## Study sessions

The four existing programs configure budgeted sessions in `course.json`.

From the repository root, preview a session:

```sh
bash course plan --course algorithms --minutes 30
```

After agreeing on the task and active minutes, start the session:

```sh
bash course start --course algorithms --minutes 30 --thread-id THREAD_ID
```

Replace `THREAD_ID` with the conversation identifier from the injected context.
Use `pause`, `resume`, and `end` for explicit session changes.
Pass the current `--thread-id` when resuming in another conversation.
The helper preserves active time across pauses and machine handoffs.
It does not send deadline alerts or automatically end sessions.

For configured untimed sessions, omit `--minutes`.
Programs without a session policy require no session commands.
See the [configuration reference](system/docs/course-configuration.md#sessions) for timing, review, and assessment policies.

For assessments, follow the selected program's `AGENTS.md` and configuration.

`context`, `plan`, and hooks do not write records.
`init` creates missing state when the learner requests saved progress.
`start` also creates missing state when study begins.
Existing progress takes precedence over initialization defaults.

## Independent study

A book or hobby program can contain sources and agent guidance without tracking anything.
Use the agent when a question, exercise, or review needs help.
There is no requirement to open Codex whenever you read or practice.

Optional state records the last reported position and work, which can become outdated after independent study.
Ask the agent to save a bookmark when useful.
Reported work does not establish assessed understanding.

## Progress across machines

Programs that save progress use `programs/<id>/state.json`.
Artifact paths in records are relative to the repository root.
Keep course work in the main checkout on `main`.
Never create or use additional worktrees or branches.

Before switching machines, pause or end any open session.
Commit and push the approved state changes on `main`.
On the other machine, pull those changes before continuing.
If state conflicts, resolve the conflict before using the helper.
Use one machine at a time for each active program.

The helper never commits, pushes, merges, or resolves conflicts.
This repository is public, so pushing state publishes its learning records.
The general learning-mode deck remains separate.

## Curriculum PDFs

Authored curricula live under `programs/<id>/curriculum/`.
Each `.tex` file is editable source, with an adjacent PDF for viewing.
The shared style is `system/latex/curriculum.sty`.

Use TeX Live 2026 with pdfTeX 3.141592653-2.6-1.40.29 and latexmk 4.88.
The build requires sourcesans 3.0 and Latin Modern fonts.
It also requires geometry, fontenc, microtype, array, enumitem, titlesec, titletoc, fancyhdr, ragged2e, multicol, url, and hyperref.

From the repository root, compile configured LaTeX sources:

```sh
make pdfs
```

The command reads the source list from each course configuration.
Books, Markdown files, and external URLs do not require PDF compilation.

Compile one curriculum by its PDF path:

```sh
make programs/algorithms/curriculum/Algorithms_Mastery_Course.pdf
```

The build replaces a viewing PDF only after compilation succeeds.
Temporary files stay under `build/pdf/`.
Use `make clean-pdfs` to remove those files while preserving viewing PDFs.
Learner PDFs under `sessions/` remain outside the curriculum build.
Render and inspect PDF changes before delivery.

## Add a program

Read the [configuration reference](system/docs/course-configuration.md).
Create a folder under `programs/` with `course.json` and `AGENTS.md`.
Use its folder name as the identifier.
Keep source paths relative to that folder and guidance paths relative to the repository root.
Keep learner exercises and their build systems inside the program.
Run `bash course validate` before using the new program.
No central course registry or Makefile course list needs updating.

## Maintenance

From the repository root, run `make build`, `make vet`, and `bash course validate`.
Exercise changed session commands with `--state /absolute/temporary/state.json` to protect actual progress.
Use `make clean-pdfs` to remove temporary PDF files after inspection.
Keep `build/bin/course`, which the shared commands and hooks require.
Preserve curricula, viewing PDFs, learner work, and progress records during cleanup.
Keep infrastructure tests absent and preserve learner-written exercise tests.
