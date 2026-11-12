# Study commands and handoff

Build the helper with `make build` from the repository root after pulling infrastructure changes.
Run the following examples from that root; replace `algorithms` with the program identifier.
Inside a program, use `bash "$(git rev-parse --show-toplevel)/course" COMMAND` and omit `--course`.

## Ordinary sessions

Read policy and saved progress, then preview the agreed active study time:

```sh
bash course context --course algorithms
bash course plan --course algorithms --minutes 30
```

After agreeing on the task and active minutes, start with the conversation identifier from the injected context:

```sh
bash course start --course algorithms --minutes 30 --thread-id THREAD_ID
```

Use `pause`, `resume`, and `end` for explicit session changes.
Pass the current `--thread-id` when resuming in another conversation.
The helper preserves active time across pauses and does not send deadline alerts or end sessions automatically.
For configured untimed sessions, omit `--minutes`. Programs without a session policy need no ordinary-session commands.

To change an ordinary session to open-ended study without closing feedback:

```sh
bash course replan --course algorithms --untimed --skip-feedback --thread-id THREAD_ID --reason "Learner requested open-ended study without closing feedback"
```

Use either flag alone when only that change is requested.
Replanning preserves the original plan and elapsed time. It cannot change assessment limits or assessment feedback.
Follow the selected program's instructions for assessments.

## Saved progress

`context`, `plan`, and hooks are read-only.
`init` creates missing state when saved progress is requested; `start` also creates it when study begins.
Existing records take precedence over initialization defaults.
Independent study can make saved position outdated. A reported bookmark does not establish assessed understanding.

Progress lives in `programs/<id>/state.json`, with artifact paths relative to the repository root.
Keep the general learning-mode deck separate.
Use the main checkout on `main`, one machine at a time for each active program.
Before switching, pause or end the session and synchronize approved Git changes.
Resolve state conflicts before using the helper.
The helper never commits, pushes, merges, or resolves conflicts.

The [configuration reference](course-configuration.md) defines timing, reviews, assessments, and initialization.
