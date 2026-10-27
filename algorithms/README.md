# Algorithms course

This folder contains the tutoring instructions and local hooks for the Algorithms Mastery course.
The PDF and its derived text are tracked together in Git.
The existing `algorithm-design-manual` repository remains separate and is excluded from this repository.
`docs/course-guide.txt` provides a text copy for reading the guide in future sessions.

## Languages

Go remains the starting language for the baseline, new concepts, and initial timed assessments.
C++ enters through selected builds and reviews, followed by fresh related tasks.
Rust starts with basic syntax, tests, collections, and rules for ownership and borrowing, then familiar algorithm exercises.

Each exercise uses one language. Selected ideas return in another language after a delay.
The tutor records algorithm feedback separately from language feedback.
Timed C++ or Rust assessments require prior evidence of independent work in that language.
Language practice fits within the chosen session duration. Assessment standards remain unchanged.
See section 24 of the guide for the language plan and section 25 for session planning.

## Start a session

Open Codex in `algorithms/` within the `courses` repository or one of its worktrees.
Tell the tutor how many active study minutes you have, for example, "I have 30 minutes."
If you omit the duration, the tutor asks before starting. There is no default session length.
If your available window includes planned breaks, subtract those breaks before calculating the study time.
Ask the tutor to start, pause, resume, or end the session.
The tutor operates the timer commands and maintains the course record.

Ordinary study reserves one-sixth each for retrieval and feedback, rounded to whole minutes with halves rounded up.
Each section gets at least one minute. The remaining minutes go to focused work.
The helper accepts durations of at least three minutes.
These splits are planning rules, not measured optimal learning intervals.

| Available minutes | Retrieval | Focused work | Feedback |
| --- | --- | --- | --- |
| 20 | 3 | 14 | 3 |
| 30 | 5 | 20 | 5 |
| 45 | 8 | 29 | 8 |
| 60 | 10 | 40 | 10 |
| 90 | 15 | 60 | 15 |
| 120 | 20 | 80 | 20 |

The tutor selects useful active reviews that fit the retrieval budget, or a baseline prompt if none are due.
A due date does not require the original exercise. Reviews use varied prompts and bounded repairs.
The sum-of-cubes identity is retired from routine review. Its history remains available.
Focused work follows the weekday theme, prerequisites, unfinished work, and recorded learner feedback.
A short session can cover one model, example, proof step, or repair.
A longer session can include implementation, tests, and a related problem, depending on progress.
For 90-120 minutes, the tutor divides focused work into smaller tasks with stopping points.
The tutor plans a demonstration, supported practice, and a different task for independent application.
The tutor explains how to choose the approach and supporting claims before expecting independent construction.
Short sessions can cover part of this sequence and retain an exact next step.
Before work starts, the tutor states the purpose, budgets, independent output, success criteria, and stopping rule.
Questions diagnose specific gaps. Sufficient reasoning does not need preferred wording.
After two unsuccessful prompts on one gap, the tutor explains directly or changes the example.
This is a teaching safeguard, not a research-derived threshold.
Feedback assesses the whole argument, code, and tests, then identifies the first substantive gap for repair.
After two teaching sessions on a targeted repair, the tutor compares independent output and help required.
If there is no improvement, the tutor changes the instruction or task scope.
Proof reviews do not indefinitely replace the unfamiliar Easy and Medium baselines or later problem-solving.
Unfinished ordinary work continues in a later session.

Timed assessments keep their fixed limits and use a separate schedule.
Reserve 60 minutes for a 45-minute benchmark and review, or 75 minutes for a 60-minute stage task and review.
A 90-minute stage task needs 105 active minutes including review, and can span planned study sessions.
The tutor agrees on each partial block and pauses the attempt at that boundary.
The assessment timer shows the whole attempt and review, not each partial block.
Do not study the solution or receive feedback before the attempt ends.
If an uninterrupted assessment and review do not fit, use ordinary study or schedule the assessment for another session.
Extra available time can fund separate work after assessment feedback. It never extends the assessment deadline.
The limits and passing standards in sections 22-23 stay unchanged.
Silence counts as work until you report a break.

## Hook setup

`.codex/hooks.json` registers `SessionStart` and `UserPromptSubmit`.
Both events call `scripts/course hook` and inject the time, weekday, progress, active due reviews, and session status.
Ordinary session context also includes the standing teaching contract. Assessment context omits teaching guidance and next-task hints.
`SessionStart` also restores this context after compaction.
Opening a chat does not start a study timer.

Open `/hooks` in the Codex CLI and review the two project hook definitions.
Trust the project and the hook definitions through Codex before relying on automatic injection.
Start or reopen a session after setup, then send a prompt to test the hook.
Ask the tutor to report the injected timestamp without a clock tool call.

The scripts use the installed `/opt/homebrew/bin/uv` and Python 3.14.8, pinned in `.python-version`.
The command runs offline and uses only the Python standard library.
The hook does not alter the global learning hooks or the chezmoi configuration.

The timestamp describes the start of a turn. It is not a continuously updated clock.
The hook does not issue alerts while you work silently.
During long turns, the tutor can refresh the clock through the helper.

Sources: [Codex hooks](https://learn.chatgpt.com/docs/hooks) and [project instructions](https://learn.chatgpt.com/docs/agent-configuration/agents-md).
Local command hooks require local orchestration. See the hook source for supported environments.

## Course record

The shared runtime record is `~/.local/state/learn/courses/algorithms/state.json`.
All worktrees use this path, so moving between them does not reset progress or timing.
The helper writes the timer. The tutor updates progress, reviews, and feedback in the same record.
The hook only reads it.

`progress` contains the stage, unit, default language, teaching contract, lesson plan, next task, algorithm evidence, and language evidence.
`reviews` contains retrieval prompts, review dates, sources, and attempt results.
`session` contains the open session, its language, owner thread, and reported breaks.
Ordinary sessions also contain `plan_minutes`, the selected split in minutes.
`history` contains ended sessions and the feedback that the tutor adds.
The timer preserves the planned split and actual elapsed time in history.
Older ordinary sessions without a stored split retain their original one-hour timing.

A review starts with this shape. The dates below are examples.

```json
{
  "id": "binary-search-invariant",
  "skill": "Binary search interval reasoning",
  "prompt": "State the search interval rule and prove that each iteration preserves it.",
  "source": "Course Stage 4",
  "language": "Go",
  "learned_on": "2026-10-05",
  "due": "2026-10-07",
  "attempts": []
}
```

Review dates start around 2, 7, 21, and 42 days after learning and adjust to recall.
The optional review `status` defaults to `active`. A `retired` prompt stays in history but is excluded from the due queue.
Retiring a prompt does not establish mastery of its skill. A replacement review can test that skill with a different task.
The tutor selects tasks from the guide, progress evidence, and the weekday pattern.
The clock does not choose tasks or establish mastery.

## Manual commands

From this repository or a worktree, read the context with:

```sh
bash scripts/course context
```

Preview a 30-minute split without starting a timer:

```sh
bash scripts/course plan --minutes 30
```

Start ordinary study with the chosen duration:

```sh
bash scripts/course start --minutes 30 --language Go --thread-id THREAD_ID
```

Replace `THREAD_ID` with the thread identifier from the injected course context.
The helper also accepts `init`, `pause`, `resume`, and `end`.
Use `start --kind assessment` for a timed assessment.
Do not pass `--minutes` for assessments. The tutor first makes sure that the available time fits the assessment schedule.
For ordinary practice, add `--language 'C++'` or `--language Rust` to the start command. The default is Go.
Add `--assessment-minutes 60` or `90` for a stage task with that limit.
For an assessment pause, supply `--infrastructure-failure --reason 'DESCRIPTION'` only for a genuine infrastructure failure.
A 90-minute task can span study sessions. Use `pause --session-boundary --reason 'End of planned study block'` for that break.
Do not study the solution between those sessions.
`init` creates missing state and never resets an existing record.

## Update the guide

Edit the PDF in a worktree. Preserve its existing exercises, links, and assessment standards.
Update every maintained derived copy and dependent instruction in the same change.
Review this README, `AGENTS.md`, helper defaults, and the live course state for related changes.
Preserve historical session evidence.

From the worktree root, regenerate the derived text:

```sh
UV_CACHE_DIR=/private/tmp/algorithms-course-uv-cache uv run --script scripts/sync_course_guide.py
```

The helper uses Python 3.14.8 and PyMuPDF 1.26.7, pinned in its script metadata.
It records the PDF's SHA256 fingerprint and extracts the text and external link destinations.
Render and inspect the changed PDF pages before delivery.
Run the same command with `--check` to make sure that the text matches the PDF.

## Tests

From the repository root, run:

```sh
UV_CACHE_DIR=/private/tmp/algorithms-course-uv-cache uv run --offline --no-project --python 3.14.8 python -m unittest discover -s tests -v
```

The tests cover variable splits, phase boundaries, required durations, historical sessions, and planning without starting a timer.
They cover teaching-contract injection, retired reviews, and the exclusion of teaching guidance during assessments.
They also cover Chicago date boundaries, daylight-saving changes, pauses, assessment limits, due reviews, and read-only hook output.
They use temporary state and do not alter the real course record.
