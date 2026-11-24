# Course configuration

Each `programs/<id>/course.json` defines stable policy for one learning program.
The shared command validates the configuration before reading or writing session state.
The decoder rejects unknown fields and unsupported format versions.

The configuration defines defaults and permitted behavior.
`AGENTS.md` supplies subject-specific teaching instructions.
An optional `state.json` records reported progress, observed learning, assistance, and any tracked sessions.
Changing defaults does not reset an existing state file.
Learning outside Codex can make recorded progress outdated.

Programs need no state file to supply sources and agent guidance.
For independent study, create or update progress only at the learner's request.
Record reported work without inventing sessions, elapsed time, or assessment evidence.

## Required fields

- `schema_version` is `1`.
- `id` matches the program folder name.
- `name` is a nonempty display name.
- `sources` contains at least one source, with the primary source first.
- `learning` defines the driver, agent role, and assistance limits.

## Sources

Each source requires exactly one `path` or `url`, plus a `format`.
Allowed formats are `latex`, `markdown`, `pdf`, `text`, and `web`.

Local paths are relative to the program folder and must identify existing files.
URLs use HTTP or HTTPS.
The helper validates URL structure without fetching remote content.
A web source requires a URL.
A LaTeX source requires a local `.tex` file.

`context` resolves local source paths and preserves external URLs.
`pdf-sources` lists only configured LaTeX files, relative to the repository root.
`make pdfs` uses that list.
A viewing PDF stays beside its LaTeX source and needs no separate source entry.

## Learning policy

`learning.driver` describes what determines the next activity:

- `curriculum` follows the program's stages, prerequisites, and published standards.
- `book` follows the selected book and the learner's current position.
- `project` follows the selected project and its milestones.
- `learner` leaves the next activity to the learner.

`learning.agent_role` describes the default interaction:

- `tutor` guides practice and protects learner decisions.
- `collaborator` works with the learner on the selected task.
- `consultant` answers requested questions without assigning work by default.

`learning.explanations` and `learning.implementation` each select one permission level:

- `on_request` requires a request or an agreed teaching or implementation block.
- `proactive` permits assistance within the authorized task.
- `never` prohibits that assistance under the current agreement.

The driver and role are independent.
A book can remain the primary source while the agent acts as a tutor or consultant.
A project can reserve implementation for the learner or permit collaboration.
A working artifact alone does not establish independent understanding.

Optional `learning.guidance` lists shared Markdown files relative to the repository root.
Every listed file must exist.
The current programs select `system/guidance/tutor.md`.
Its [closing feedback process](../guidance/tutor.md#closing-feedback) applies to every program that selects it.
Keep general feedback rules in that shared guidance and subject-specific criteria in each program.
Course context supplies the selected policy and guidance paths to the agent.
These fields are instructions for the agent, not file-access controls.
The helper enforces configuration validity, session budgets, and assessment pause rules.

## Sessions

Omitting `session` disables ordinary session tracking.
Ordinary `plan` and `start` commands then fail without writing state.
Context and hooks still supply sources, guidance, and any existing progress.
The learner can ask questions without starting a session or creating state.

An explicit untimed policy, `{"timing": "untimed"}`, enables elapsed-time tracking.
Untimed study rejects `--minutes` and records elapsed time without a countdown.
Its `plan` result is an empty object.

A budgeted policy requires `timing: "budgeted"` and a `split` object.
The split contains integer weights for `retrieval` and `focused_work`, plus a fixed `feedback_minutes` budget.
The focused-work weight must be positive.
A zero retrieval weight or a zero feedback budget disables that phase.

The helper uses the two weights to calculate the share of active minutes for retrieval.
It rounds retrieval to the nearest minute, with halves rounded up.
If retrieval is enabled, it receives at least one minute.
The helper reserves the fixed feedback budget and gives the remaining minutes to focused work.
If the plan leaves no time for focused work, the command fails before it changes state.

The current programs use weights `1` and `5` and reserve three minutes for ordinary feedback.
This gives retrieval one-sixth of the total time before rounding.
With these settings, a session needs at least five active minutes.
No session starts from a hook, a context read, or a plan preview.
Changes to policy do not rewrite a session's saved budgets.

An explicit learner request can override the current ordinary session without changing course defaults.
Use `replan --untimed` to remove its countdown or `replan --skip-feedback` to omit its closing feedback block.
Both flags can be combined. Supply the current session owner's `--thread-id` and a `--reason` describing the request.
The command rejects assessments, missing sessions, and minute-limit arguments.

The helper stores the effective change in `session.study_override` with its time and reason.
It preserves the original `plan_minutes`, start time, and pauses.
Context and hooks expose `study_override` and apply it to the current session summary.
An untimed summary has elapsed time but no phase or remaining-time fields.
Skipping feedback alone removes that phase from the effective plan.
The override remains with the session when `end` moves it into history.
Keep the lesson plan aligned with this effective state. Archive superseded plans instead of retaining conflicting current instructions.

## Reviews

Optional `reviews` contains `interval_days` and `weekdays_only`.
Intervals must form a nonempty, strictly increasing list of positive integers.
The current programs use `[2, 7, 21, 42]` and move weekend reviews to the next study day.

The agent uses this policy when recording reviews and explains any adjustment.
The helper reports due, nonretired reviews from state.
It does not generate reviews, schedule dates, or judge mastery.
Without a review policy, context omits due-review counts and identifiers.
Existing review records remain preserved.

## Assessments

Optional `assessment` contains `minutes`, `feedback`, and `split_minutes`.
Allowed durations must be positive and strictly increasing.
Feedback requires at least one minute.
Each split duration must appear in the allowed duration list.

Algorithms and Algorithms Fun define timed assessments of 45, 60, or 90 minutes, plus 15 feedback minutes.
In both programs, only a 90-minute assessment permits an agreed session boundary.
Infrastructure failures remain a separate allowed pause reason.
Ordinary session timing does not change assessment limits.

Assessment context hides ordinary teaching policy, source lists, task plans, and teaching feedback.
Its primary curriculum reference and timing remain available.
The agent must preserve unaided attempts, including when the configuration permits implementation during ordinary work.

## Initial progress

Optional `initial` seeds missing progress when `init` or `start` creates state.
Allowed fields are `stage`, `language`, `unit`, `position`, and `next_task`.
Stage is a nonnegative integer.
The remaining fields are text.
Use `position` for a chapter, section, or milestone when a stage number does not fit.

The helper creates empty evidence, review, and history collections.
It does not create completed work from configuration.
Context and hooks do not present these defaults as recorded progress before initialization.
An existing state file supplies recorded progress and language.
Treat a saved position as the last reported position, not proof of where the learner is now.

Optional `weekday_work` contains seven descriptions, starting with Monday.
It appears in ordinary context but does not schedule sessions.

## Book example

This example requires no sessions, saved progress, or scheduled reviews.
It follows the book, permits requested explanations, and reserves implementation for the learner.

```json
{
  "schema_version": 1,
  "id": "rust-book",
  "name": "The Rust Programming Language",
  "sources": [
    {
      "url": "https://doc.rust-lang.org/book/",
      "format": "web"
    }
  ],
  "learning": {
    "driver": "book",
    "agent_role": "consultant",
    "explanations": "on_request",
    "implementation": "never"
  }
}
```

Use observed progress when migrating an existing program.

## Validation

From the repository root, validate all configurations:

```sh
bash course validate
```

To validate one program, add `--course ID`.
The command reports validated programs and exits with an error for invalid configuration.
Session commands use the same loader before any state write.
