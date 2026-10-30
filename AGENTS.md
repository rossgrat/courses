# Learning programs

## Scope and sources

- Apply the selected program's learning policy during study. Planning, maintenance, and general questions do not start a session.
- Read `programs/<id>/course.json`, its `AGENTS.md`, and its existing `state.json` before teaching.
- Read the files listed in `learning.guidance` before ordinary study. Resolve those paths from the repository root.
- Use the configured sources for the current task. Resolve source paths from the program folder. Keep external sources as URLs.
- Follow the configured driver, agent role, and assistance limits. Keep casual questions within the requested scope.
- Treat the configuration as stable policy. Treat state as the saved learning record. Use `initial` only when creating missing state.
- A program needs no state file or tracked sessions. Do not require either for reading, building, or asking questions.
- Treat saved progress as the last reported position. Learning outside this conversation can make that position outdated.
- For independent study, save progress only when requested. Record reported work without inventing sessions, elapsed time, or assessment evidence.
- Keep subject rules, prerequisites, languages, and assessment standards in the program folder.
- Follow more specific instructions when a program or exercise adds an `AGENTS.md`.
- Keep learner code, tests, notes, and workspace setup learner-owned. Change them only within the user's explicit authorization and assistance policy.
- Keep course records separate from the general learning-mode deck. Switch learning mode only when requested.
- Do all course work in the main checkout on `main`, including timer and progress records.
- Never create or use additional worktrees or branches. This repository rule overrides general worktree guidance.
- Do not add infrastructure tests. Preserve learner-written exercise tests.

## Shared commands

- Run `bash "$(git rev-parse --show-toplevel)/course" COMMAND` from a program folder.
- From the repository root, add `--course ID` to select a program explicitly.
- Use `list` to find programs. Use `validate` to read and validate every course configuration without writing records.
- If `build/bin/course` is missing, run `make build` from the repository root before relying on hooks.
- Use the exact Go version in `system/cli/.go-version`. Keep dependencies pinned to exact versions.
- Use `context` to read policy and existing records. Use `plan` only when ordinary sessions are configured.
- If `session` is omitted, do not start ordinary sessions. Explicit `untimed` sessions still record elapsed time.
- For budgeted sessions, pass `--minutes MINUTES`. For untimed sessions, omit `--minutes`.
- Treat `list`, `validate`, `pdf-sources`, hooks, `context`, and `plan` as read-only.
- Use `init` only when the user requests saved progress. It creates missing state without replacing existing progress.
- Use `make pdfs` from the repository root to compile the configured LaTeX sources.
- Preserve curriculum content, exercises, links, and historical evidence during infrastructure changes.
- Update affected instructions, course configuration, and documentation together when behavior changes.

## Time and sessions

- Apply session procedures only to configured sessions and assessments. Answer casual questions without requiring a session.
- Use `America/Chicago` for course dates. Treat injected time as the time at turn start.
- Refresh with `context` after a long turn. If injected context is missing, run `context` and report the missing injection.
- Obtain active study minutes before starting a budgeted session. Exclude planned breaks and never assume an hour.
- Use the configured session split. Reserve the displayed retrieval and feedback budgets when those phases are enabled.
- State the objective, independent output, success criteria, and stopping point before starting. Include budgets for timed work.
- For a configured session, run `start --thread-id THREAD_ID` immediately before its first task. Add `--minutes MINUTES` for budgeted study.
- Translate explicit pause, resume, and end requests into `pause`, `resume`, and `end`.
- Pass the current `--thread-id` when starting or resuming. Resuming assigns the session to that conversation.
- Silence counts as working time. Do not infer a break or reset a timer from inactivity.
- Resolve a session from another conversation or day with the user before changing it.
- Keep assessment standards fixed. Use only the limits and pause rules defined by the selected program.
- Start assessments immediately before revealing the task. Count solving, implementation, explanation, and resubmissions within the limit.
- Keep assessments unaided. Do not reveal hints, topic labels, templates, or solution code during the attempt.
- During assessments, do not expose ordinary guidance, saved task plans, initial tasks, or teaching feedback.
- Do not start an assessment when the available time cannot support its required attempt and feedback.
- Hooks do not send deadline alerts or automatically end sessions.

## Records and handoff

- Read `programs/<id>/state.json` immediately before every update.
- Preserve `progress`, `reviews`, `session`, `history`, and existing subject-specific evidence.
- Read recorded teaching feedback and lesson plans before selecting the next task.
- Let the helper manage timer fields. Record learning evidence and feedback without rewriting completed timing.
- When reviews are configured, record IDs, skills, prompts, sources, languages, learned dates, due dates, attempts, and retirement reasons.
- Use the configured review intervals and weekday policy. Record the reason for adjustments.
- Distinguish independent success, assisted practice, and misses. Record observed errors and assistance.
- Do not infer mastery from generated work, completed reading, prompted agreement, or the existence of an artifact.
- At session end, run `end` and complete the new history entry with tasks, artifacts, feedback, and unfinished work.
- Preserve the planned split and elapsed time. Set an exact next task and stopping point.
- Store artifact paths relative to the repository root. Keep machine-local paths out of new progress records.
- Before switching machines, pause or end any open session. Save existing records through Git on `main` with authorization.
- Before continuing elsewhere, synchronize committed state and resolve any state conflict before using the helper.
- Do not edit the same program state concurrently on different machines. Never silently choose one side of a conflict.
- Do not automatically commit, push, or merge progress. Obtain authorization before publishing new records.
- Never auto-grade, advance stages, or alter records from a hook.
