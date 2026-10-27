import argparse
import json
import shlex
import subprocess
import sys
import tempfile
import unittest
from datetime import datetime, timedelta, timezone
from pathlib import Path

from scripts import course


def instant(value):
    return datetime.fromisoformat(value).astimezone(timezone.utc)


def action(command):
    return argparse.Namespace(command=command, kind="study", thread_id="test-thread",
                              language="Go", minutes=60,
                              assessment_minutes=45, infrastructure_failure=False,
                              session_boundary=False, reason="Break")


class CourseTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.path = Path(self.temp.name) / "state.json"
        self.state = course.read_state(self.path)

    def snapshot(self, when):
        return json.loads(course.context(self.state, instant(when), self.path, "test-thread").split("\n", 1)[1])

    def start(self, when="2026-10-06T00:00:00+00:00", kind="study"):
        args = action("start")
        args.kind = kind
        if kind == "assessment":
            args.minutes = None
        course.change_session(self.state, args, instant(when))

    def test_chicago_day_and_due_reviews_before_local_midnight(self):
        self.state["reviews"] = [{"id": "due", "due": "2026-10-05"},
                                 {"id": "tomorrow", "due": "2026-10-06"}]
        snapshot = self.snapshot("2026-10-06T00:30:00+00:00")
        self.assertEqual(snapshot["weekday"], "Monday")
        self.assertEqual(snapshot["local_time"], "2026-10-05T19:30:00-05:00")
        self.assertEqual(snapshot["due_review_ids"], ["due"])
        self.assertEqual(snapshot["session"], "inactive")
        self.assertFalse(self.path.exists())

    def test_pause_resume_and_end_preserve_active_time(self):
        self.start()
        course.change_session(self.state, action("pause"), instant("2026-10-06T00:10:00+00:00"))
        self.assertEqual(self.snapshot("2026-10-06T00:20:00+00:00")["session"]["active_seconds"], 600)
        course.change_session(self.state, action("resume"), instant("2026-10-06T00:20:00+00:00"))
        course.change_session(self.state, action("end"), instant("2026-10-06T00:30:00+00:00"))
        self.assertIsNone(self.state["session"])
        self.assertEqual(self.state["history"][0]["active_seconds"], 1200)

    def test_retired_reviews_stay_in_history_and_leave_due_queue(self):
        self.state["reviews"] = [
            {"id": "old-proof", "due": "2026-10-01", "status": "retired",
             "attempts": [{"date": "2026-10-01", "result": "assisted"}]},
            {"id": "new-task", "due": "2026-10-05"},
        ]
        original = json.dumps(self.state)
        snapshot = self.snapshot("2026-10-06T00:30:00+00:00")
        self.assertEqual(snapshot["due_review_ids"], ["new-task"])
        self.assertEqual(snapshot["reviews_due"], 1)
        self.assertEqual(json.dumps(self.state), original)

    def test_legacy_state_receives_teaching_contract_without_mutation(self):
        del self.state["progress"]["teaching_contract"]
        original = json.dumps(self.state)
        snapshot = self.snapshot("2026-10-06T00:30:00+00:00")
        self.assertEqual(snapshot["teaching_contract"], course.TEACHING_CONTRACT)
        self.assertEqual(json.dumps(self.state), original)

    def test_ending_while_paused_does_not_count_break(self):
        self.start()
        course.change_session(self.state, action("pause"), instant("2026-10-06T00:10:00+00:00"))
        course.change_session(self.state, action("end"), instant("2026-10-06T01:00:00+00:00"))
        self.assertEqual(self.state["history"][0]["active_seconds"], 600)

    def test_daylight_saving_fallback_counts_real_elapsed_time(self):
        self.start("2026-11-01T01:30:00-05:00")
        snapshot = self.snapshot("2026-11-01T01:30:00-06:00")
        self.assertEqual(snapshot["session"]["active_seconds"], 3600)
        self.assertEqual(snapshot["session"]["phase"], "session time complete")

    def test_study_phase_boundaries(self):
        self.start()
        for time, phase in [("00:09:59", "retrieval"), ("00:10:00", "focused work"),
                            ("00:50:00", "feedback"), ("01:00:00", "session time complete")]:
            with self.subTest(time=time):
                self.assertEqual(self.snapshot(f"2026-10-06T{time}+00:00")["session"]["phase"], phase)

    def test_variable_study_plans(self):
        for minutes, parts in [(20, (3, 14, 3)), (30, (5, 20, 5)), (45, (8, 29, 8)),
                               (60, (10, 40, 10)), (90, (15, 60, 15)), (120, (20, 80, 20))]:
            with self.subTest(minutes=minutes):
                plan = course.study_plan(minutes)
                self.assertEqual(tuple(plan.values()), parts)
                self.assertEqual(sum(plan.values()), minutes)
        for minutes in [None, -1, 0, 1, 2]:
            with self.subTest(minutes=minutes), self.assertRaises(ValueError):
                course.study_plan(minutes)

    def test_start_requires_explicit_study_duration(self):
        args = action("start")
        args.minutes = None
        with self.assertRaisesRegex(ValueError, "minutes"):
            course.change_session(self.state, args, instant("2026-10-06T00:00:00+00:00"))
        self.assertIsNone(self.state["session"])

    def test_short_session_boundaries_and_saved_plan(self):
        args = action("start")
        args.minutes = 20
        course.change_session(self.state, args, instant("2026-10-06T00:00:00+00:00"))
        for time, phase, remaining in [("00:02:59", "retrieval", 1021),
                                        ("00:03:00", "focused work", 1020),
                                        ("00:17:00", "feedback", 180),
                                        ("00:20:00", "session time complete", 0)]:
            with self.subTest(time=time):
                session = self.snapshot(f"2026-10-06T{time}+00:00")["session"]
                self.assertEqual(session["phase"], phase)
                self.assertEqual(session["remaining_seconds"], remaining)
                self.assertEqual(session["planned_minutes"], 20)
        course.change_session(self.state, action("end"), instant("2026-10-06T00:20:00+00:00"))
        self.assertEqual(self.state["history"][0]["plan_minutes"],
                         {"retrieval": 3, "focused work": 14, "feedback": 3})

    def test_legacy_session_keeps_original_hour_without_rewriting_state(self):
        self.start()
        del self.state["session"]["plan_minutes"]
        original = json.dumps(self.state)
        self.assertEqual(self.snapshot("2026-10-06T00:50:00+00:00")["session"]["phase"], "feedback")
        self.assertEqual(json.dumps(self.state), original)

    def test_plan_command_does_not_start_session_or_write_state(self):
        course.save_state(self.path, self.state)
        original = self.path.read_bytes()
        result = subprocess.run([sys.executable, "scripts/course.py", "plan", "--minutes", "30",
                                 "--state", str(self.path)],
                                text=True, capture_output=True, check=True)
        self.assertEqual(json.loads(result.stdout),
                         {"retrieval": 5, "focused work": 20, "feedback": 5})
        self.assertEqual(self.path.read_bytes(), original)

    def test_study_minutes_cannot_change_assessment_limit(self):
        args = action("start")
        args.kind = "assessment"
        args.minutes = 20
        with self.assertRaisesRegex(ValueError, "assessment"):
            course.change_session(self.state, args, instant("2026-10-06T00:00:00+00:00"))
        self.assertIsNone(self.state["session"])

    def test_missing_duration_does_not_create_state(self):
        result = subprocess.run([sys.executable, "scripts/course.py", "start", "--state", str(self.path)],
                                text=True, capture_output=True)
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("--minutes", result.stderr)
        self.assertFalse(self.path.exists())

    def test_all_assessment_limits_keep_fifteen_minutes_for_feedback(self):
        start = instant("2026-10-06T00:00:00+00:00")
        for minutes in [45, 60, 90]:
            with self.subTest(minutes=minutes):
                args = action("start")
                args.kind = "assessment"
                args.minutes = None
                args.assessment_minutes = minutes
                course.change_session(self.state, args, start)
                deadline = start + timedelta(minutes=minutes)
                session = self.snapshot(deadline.isoformat())["session"]
                self.assertEqual(session["phase"], "feedback")
                self.assertEqual(session["remaining_seconds"], 900)
                self.assertEqual(session["assessment_remaining_seconds"], 0)
                course.change_session(self.state, action("end"), deadline + timedelta(minutes=15))

    def test_assessment_deadline_and_pause_restriction(self):
        self.start(kind="assessment")
        with self.assertRaisesRegex(ValueError, "infrastructure failure"):
            course.change_session(self.state, action("pause"), instant("2026-10-06T00:10:00+00:00"))
        before = self.snapshot("2026-10-06T00:44:59+00:00")
        after = self.snapshot("2026-10-06T00:45:00+00:00")
        self.assertEqual(before["session"]["phase"], "assessment")
        self.assertEqual(after["session"]["phase"], "feedback")
        self.assertNotIn("next_task", after)
        self.assertNotIn("unit", after)
        self.assertNotIn("teaching_contract", before)
        self.assertNotIn("teaching_contract", after)

    def test_ninety_minute_task_can_span_days_without_counting_break(self):
        args = action("start")
        args.kind = "assessment"
        args.assessment_minutes = 90
        args.minutes = None
        course.change_session(self.state, args, instant("2026-10-06T00:00:00+00:00"))
        pause = action("pause")
        pause.session_boundary = True
        course.change_session(self.state, pause, instant("2026-10-06T00:50:00+00:00"))
        course.change_session(self.state, action("resume"), instant("2026-10-07T00:00:00+00:00"))
        snapshot = self.snapshot("2026-10-07T00:39:00+00:00")
        self.assertEqual(snapshot["session"]["active_seconds"], 89 * 60)
        self.assertEqual(snapshot["session"]["assessment_remaining_seconds"], 60)
        self.assertEqual(snapshot["session"]["phase"], "assessment")

    def test_existing_session_cannot_be_reset_by_start(self):
        self.start()
        with self.assertRaisesRegex(ValueError, "already open"):
            self.start("2026-10-07T00:00:00+00:00")
        self.assertEqual(self.state["session"]["started_at"], "2026-10-06T00:00:00+00:00")

    def test_hook_both_events_reads_saved_state_without_writing(self):
        self.state["progress"]["teaching_contract"] = "Demonstrate, practice, and assess a whole argument."
        course.save_state(self.path, self.state)
        original = self.path.read_bytes()
        hooks = json.loads(Path(".codex/hooks.json").read_text())["hooks"]
        for event in ["SessionStart", "UserPromptSubmit"]:
            with self.subTest(event=event):
                command = hooks[event][0]["hooks"][0]["command"]
                result = subprocess.run(
                    command + " --state " + shlex.quote(str(self.path)),
                    input=json.dumps({"hook_event_name": event, "session_id": "test-thread"}),
                    cwd="tests", shell=True, executable="/bin/bash",
                    text=True, capture_output=True, check=True,
                )
                output = json.loads(result.stdout)["hookSpecificOutput"]
                self.assertEqual(output["hookEventName"], event)
                self.assertIn("America/Chicago", output["additionalContext"])
                self.assertIn(self.state["progress"]["teaching_contract"], output["additionalContext"])
                self.assertIn('"thread_id": "test-thread"', output["additionalContext"])
                self.assertEqual(self.path.read_bytes(), original)

    def test_init_does_not_reset_existing_progress(self):
        self.state["progress"]["stage"] = 4
        course.save_state(self.path, self.state)
        subprocess.run([sys.executable, "scripts/course.py", "init", "--state", str(self.path)],
                       text=True, capture_output=True, check=True)
        self.assertEqual(course.read_state(self.path)["progress"]["stage"], 4)

    def test_language_selection_survives_session_history(self):
        for language in ["Go", "C++", "Rust"]:
            with self.subTest(language=language):
                args = action("start")
                args.language = language
                course.change_session(self.state, args, instant("2026-10-06T00:00:00+00:00"))
                self.assertEqual(self.snapshot("2026-10-06T00:05:00+00:00")["language"], language)
                course.change_session(self.state, action("end"), instant("2026-10-06T01:00:00+00:00"))
                self.assertEqual(self.state["history"][-1]["language"], language)


if __name__ == "__main__":
    unittest.main()
