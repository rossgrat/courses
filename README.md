# Courses

Personal learning programs with their own sources, teaching rules, and optional progress records:

- [Algorithms](programs/algorithms/)
- [Algorithms Fun](programs/algorithms-fun/)
- [Trading Engineering](programs/trading-engineer/)
- [Oxide Systems Engineering](programs/oxide/)
- [Music Technology Engineering](programs/music-technology/)

## Start

Use Go 1.27.3, as pinned in [system/cli/.go-version](system/cli/.go-version).
From the repository root:

```sh
make build
bash course list
bash course validate
bash course context --course algorithms
```

Open the selected `programs/<id>/` folder for study.
Hooks select a course from the working directory; the repository root and separate nested Git repositories do not select one.
Independent reading and practice do not require an agent or a tracked session.

## Sessions and records

Keep all work in the main checkout on `main`. Do not create branches or worktrees.
This repository is public; pushing progress publishes learning records.
Pause or end a session before switching machines, then synchronize approved changes and resolve conflicts before continuing.

See [study commands and handoff](system/docs/study.md), [curriculum builds](system/docs/curriculum-builds.md), and the [configuration reference](system/docs/course-configuration.md).
The root [AGENTS.md](AGENTS.md) and each program's instructions define ownership and assessment rules.
