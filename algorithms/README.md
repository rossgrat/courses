# Algorithms course

`Algorithms_Mastery_Course.pdf` defines the curriculum.
`docs/course-guide.txt` contains its derived text and links.
`AGENTS.md` adds algorithm-specific teaching and assessment rules to the shared parent instructions.
`state.json` stores current progress, reviews, and session history in Git.

Follow the [shared setup and handoff instructions](../README.md) before starting on a new machine.
Open Codex in this folder and give the tutor your available active study minutes.
Go is the starting language. Section 24 defines the C++ and Rust progression.

## Commands

From this folder, use the compatibility launcher:

```sh
bash scripts/course context
bash scripts/course plan --minutes 30
bash scripts/course start --minutes 30 --language Go --thread-id THREAD_ID
```

Replace `THREAD_ID` with the conversation identifier in the injected context.
Use `pause`, `resume`, and `end` for explicit session changes.
Pass the current `--thread-id` when resuming in another conversation.

## Assessments

Use `start --kind assessment` for a 45-minute benchmark plus 15 minutes of feedback.
Use `--assessment-minutes 60` or `90` for the stage tasks specified in section 22.
Do not pass `--minutes` for assessments.
The assessment limit includes implementation, explanation, debugging, and resubmissions.

A 90-minute stage task can span agreed sessions.
Use `pause --session-boundary --reason DESCRIPTION` at the boundary.
Give no feedback or solution study before the attempt ends.
Use `--infrastructure-failure --reason DESCRIPTION` for an infrastructure pause.

## Course files

After a curriculum change, regenerate and compare the guide from this folder:

```sh
bash scripts/course sync-guide
bash scripts/course sync-guide --check
```

Preserve the source exercises, links, and assessment standards.
The separate `algorithm-design-manual` repository remains ignored by `courses`.
Learner-written exercises and their tests remain under `sessions/`.
