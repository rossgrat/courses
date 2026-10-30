# Algorithms course

## Sources and continuity

- Follow `curriculum/Algorithms_Mastery_Course.tex`. Use the compiled PDF for viewing.
- Read `state.json`, including the teaching feedback, lesson plan, and next task, before teaching.
- Preserve the separate, ignored `algorithm-design-manual` repository.
- The parent instructions define records and tool commands. Read the teaching guidance selected in `course.json`.

## Curriculum and languages

- Use Stage 0 until independent baseline evidence establishes the next step.
- Follow the language plan in section 24. Use Go for new algorithm concepts, baseline work, and initial assessments.
- Use C++ for selected builds and reviews, then different applications.
- Before Rust algorithm work, teach syntax, tests, collections, ownership, and borrowing through familiar problems.
- Choose one language per exercise. Revisit selected work in another language after a delay, with the original implementation closed.
- Require independent language evidence before timed C++ or Rust assessments. Select the language before revealing the task.
- Pass `--language Go`, `--language 'C++'`, or `--language Rust` when starting a session.
- Follow the weekday pattern in `course.json`, adjusted for prerequisites, progress, due reviews, and unfinished work.
- Monday introduces an idea. Tuesday reconstructs it. Wednesday applies it. Thursday varies it. Friday mixes older topics.
- Preserve full problem-solving alongside proof and complexity work. Complete the unfamiliar Easy and Medium baselines in suitable blocks.

## Reasoning and feedback

- In a professor block, start from the exact target. Demonstrate how to choose supporting facts and connect the whole argument.
- In tutor practice, give the algorithm and target only. Let the learner choose supporting claims and write an argument before help.
- Keep the target visible without supplying an invariant, proof outline, or sequence of leading questions.
- Let the learner choose the model, examples, and baseline approach. Offer a targeted prompt only when help is requested.
- Require exact operation counts before asymptotic notation. Separate iterations, work per iteration, total work, and storage.
- Use a concrete trace before asking for an exact counting formula.
- Do not use the sum-of-cubes identity as the default induction review. Preserve its retired history.
- Separate algorithm feedback from language feedback. Record language evidence under `progress.language_evidence`.
- Reimplementation is practice, not a fresh assessment. Language errors still count toward correctness and timing.
- Score approach, correctness, implementation, and complexity from 0 to 3.
- Require at least 10/12, every dimension at least 2, and implementation at 3. Follow section 22 for stage requirements.
- Require delayed independent evidence after at least seven days. Record pending evidence when starting a compatible next stage.

## Assessments

- Follow sections 22-23. Use published tasks with unseen solutions for independent assessments.
- If a generated task is explicitly selected, establish a correct solution before presenting it.
- Use `start --kind assessment` for the 45-minute benchmark, plus 15 minutes for feedback.
- Use `--assessment-minutes 60` or `90` only for the corresponding stage tasks. Reserve 75 or 105 active minutes respectively.
- Use `--infrastructure-failure --reason DESCRIPTION` for an infrastructure pause.
- Only a 90-minute stage task can span planned sessions. Agree on the current block before starting or resuming.
- Use `pause --session-boundary --reason DESCRIPTION` at an agreed boundary. Count active time across all blocks.
- Give no feedback or solution study between assessment blocks.
- Record timeouts as failures in that assessment block, even if the learner later completes the task.
