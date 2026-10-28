# Courses

This repository contains shared study tools and separate curricula for Algorithms, Trading Engineering, and Oxide Systems Engineering.
The root `AGENTS.md` defines shared tutoring behavior.
Each course adds its curriculum rules through its own `AGENTS.md` and `course.json`.

## Setup

Install Go 1.27.3 on each machine.
From the repository root, build the tool:

```sh
make build
```

The build produces `bin/course`, which Git ignores.
The executable uses the Go standard library and includes time-zone data.
After pulling infrastructure changes, run `make build` again.

Open Codex in the selected course folder.
Review and trust the two shared definitions through `/hooks`.
Start a new chat after changing instructions or hook definitions.

Codex loads instructions from the Git root down to the chat's working directory.
Matching hooks from parent and child sources all run.
The shared course hook uses the session working directory to select one course.
Opening the repository root injects no course context.
Sources: [instructions](https://learn.chatgpt.com/docs/agent-configuration/agents-md) and [hooks](https://learn.chatgpt.com/docs/hooks).

## Study commands

From any course folder, read its context:

```sh
bash "$(git rev-parse --show-toplevel)/scripts/course" context
```

From the repository root, select a course explicitly:

```sh
bash scripts/course context --course algorithms
bash scripts/course context --course trading-engineer
bash scripts/course context --course oxide
```

The helper recognizes the enclosing course from deeper exercise folders.
It stops at the Git root, so a separate nested repository does not inherit course selection.

Preview a 30-minute session from the repository root:

```sh
bash scripts/course plan --minutes 30
```

The ordinary split reserves one-sixth each for retrieval and feedback, rounded to whole minutes.
The remaining minutes cover focused work.
These are planning defaults, not curriculum deadlines or measured optimal learning intervals.

After selecting the duration and task, start from the repository root:

```sh
bash scripts/course start --course algorithms --minutes 30 --thread-id THREAD_ID
```

Replace `THREAD_ID` with the conversation identifier in the injected context.
Use `pause`, `resume`, and `end` for explicit session changes.
Pass the current `--thread-id` when resuming in another conversation.
The helper preserves elapsed active time across pauses and machine handoffs.
It does not send deadline alerts.

`context`, `plan`, and hooks do not write records or start timers.
`context` and both hooks include the saved teaching contract during ordinary study.
Assessment context omits the contract and task guidance to preserve unaided attempts.
`init` creates a missing course record without starting study.
`start` also creates a missing record when study begins.
Only Algorithms defines timed assessments. Its local instructions explain the limits.

During practice, the tutor gives the task and waits for the learner to choose an approach.
The tutor answers questions narrowly and leaves claim selection and repairs to the learner.
A requested or agreed professor block explains choices through a worked example.
The tutor role resumes after that example.
Supported practice makes help available without revealing it before an attempt.

## Progress across machines

Each course stores progress, reviews, the current session, and completed history in `<course>/state.json`.
Git carries that record between machines along with the curriculum and learner artifacts.
Artifact paths in new records are relative to the repository root.
Each checkout or worktree has its own record until Git synchronizes changes.

Before switching machines, pause or end the session.
Commit the state changes and push the approved commit.
On the other machine, pull those changes before continuing.
If a state conflict occurs, resolve it before using the helper.
Use one working copy at a time for each active course.

The helper never commits, pushes, merges, or chooses a conflict resolution.
The repository is public, so pushing a state file publishes its learning records.

The migration imports the existing Algorithms record and converts its repository paths to relative paths.
The old `~/.local/state/learn/courses/algorithms/state.json` remains a backup and is no longer read by this tool.
The separate general learning-mode deck remains outside this repository.

## Curriculum PDFs

Each course's `.tex` file contains its editable curriculum, tables, and links.
The tutor reads this LaTeX source directly.
The adjacent PDF is the compiled version for viewing.
The three courses use the standard LaTeX `article` layout.

Use TeX Live 2026 with pdfTeX 3.141592653-2.6-1.40.29 and latexmk 4.88.
The build requires Latin Modern, fontenc, array, url, and hyperref for fonts, tables, and links.
From the repository root, compile all three curricula:

```sh
make pdfs
```

Compile one curriculum with its PDF path:

```sh
make algorithms/Algorithms_Mastery_Course.pdf
```

The build replaces each course PDF only after its LaTeX compilation succeeds.
Contents pages and curriculum-map page references update during compilation.
Intermediate files stay in the ignored `build/pdf/` directory.
Use `make clean-pdfs` to remove those files while preserving the viewing PDFs.
The PDFs under `algorithms/sessions/` are outside this build.
Keep related instructions and course configuration consistent with curriculum changes.
Render and inspect compiled PDFs before delivery.

## Add a course

Create its folder and add the LaTeX source.
Add a `course.json` with `id`, `name`, `source`, `language`, `unit`, and `next_task`.
Set `id` to the course folder name.
Set `source` to the `.tex` path relative to that folder.
Use an existing course file as the format reference.
Add `AGENTS.md` for subject-specific rules and exceptions to shared defaults.
Add the course folder to `TEX_SOURCES` in the root `Makefile`.

Optional `weekday_work` contains seven entries, starting with Monday.
Optional `assessment` defines allowed durations, feedback minutes, and durations that permit session boundaries.
Course-specific tools belong in that course's `scripts/` folder.
Add local hooks only for distinct tasks. Matching hooks do not override shared hooks.

## Infrastructure maintenance

Use `make build` and `make vet` from the repository root.
Exercise changed commands manually with `--state /absolute/temporary/state.json` to protect actual progress.
Keep infrastructure tests absent and preserve learner-written exercise tests.
The Algorithms `scripts/course` launcher remains as a compatibility entry point.
