import argparse
import json
import sys
from datetime import datetime, timezone
from pathlib import Path
from zoneinfo import ZoneInfo


STATE_PATH = Path.home() / ".local/state/learn/courses/algorithms/state.json"
CHICAGO = ZoneInfo("America/Chicago")
TEACHING_CONTRACT = (
    "Plan each lesson around an independent result and its place in the course. "
    "Demonstrate how to choose the approach and supporting claims, then reduce help through "
    "supported practice and a different independent task. Keep the whole argument visible. "
    "Questions diagnose gaps; accept sufficient wording. After two unsuccessful prompts on one gap, "
    "explain directly or change the example. Select useful, varied reviews within budget; do not "
    "default to the sum-of-cubes proof. Judge whole-task performance and delayed application, "
    "not prompted agreement. After two repair sessions, compare independent output and help required. "
    "Change instruction if progress stalls. Keep fresh problem-solving in the plan."
)
DAY_WORK = (
    "Study one small idea and explain a worked example",
    "Reconstruct, implement, and test without the resource",
    "Solve a related problem with help only when needed",
    "Change a constraint or solve another application",
    "Review older material and solve mixed, unlabeled problems",
    "No scheduled session",
    "No scheduled session",
)


def read_state(path):
    if path.exists():
        return json.loads(path.read_text())
    return {
        "progress": {
            "stage": 0,
            "unit": "Baseline and setup",
            "language": "Go",
            "teaching_contract": TEACHING_CONTRACT,
            "next_task": "Attempt selected MIT 6.006 Problem Set 0 questions before solutions",
            "evidence": [],
            "language_evidence": {"Go": [], "C++": [], "Rust": []},
        },
        "reviews": [],
        "session": None,
        "history": [],
    }


def save_state(path, state):
    path.parent.mkdir(parents=True, exist_ok=True)
    temporary = path.with_suffix(".tmp")
    temporary.write_text(json.dumps(state, indent=2) + "\n")
    temporary.replace(path)


def seconds_between(start, end):
    return (datetime.fromisoformat(end) - datetime.fromisoformat(start)).total_seconds()


def elapsed_seconds(session, now):
    end = session["ended_at"] or now.isoformat()
    breaks = sum(seconds_between(p["start"], p["end"] or end) for p in session["pauses"])
    return int(seconds_between(session["started_at"], end) - breaks)


def is_paused(session):
    return bool(session["pauses"] and session["pauses"][-1]["end"] is None)


def study_plan(minutes):
    if minutes is None:
        raise ValueError("Give the available study time with --minutes before starting.")
    if minutes < 3:
        raise ValueError("Use at least 3 study minutes to leave time for each section.")
    review = max(1, (minutes + 3) // 6)
    return {"retrieval": review, "focused work": minutes - 2 * review, "feedback": review}


def context(state, now, path, thread_id):
    local = now.astimezone(CHICAGO)
    due = [r["id"] for r in state["reviews"]
           if r.get("status", "active") != "retired" and r["due"] <= local.date().isoformat()]
    result = {
        "course": "Algorithms Mastery",
        "local_time": local.isoformat(timespec="seconds"),
        "timezone": "America/Chicago",
        "weekday": local.strftime("%A"),
        "weekday_work": DAY_WORK[local.weekday()],
        "thread_id": thread_id,
        "state_file": str(path),
        "stage": state["progress"]["stage"],
        "unit": state["progress"]["unit"],
        "language": state["progress"].get("language", "Go"),
        "next_task": state["progress"]["next_task"],
        "teaching_contract": state["progress"].get("teaching_contract", TEACHING_CONTRACT),
        "reviews_due": len(due),
        "due_review_ids": due[:10],
        "session": "inactive",
    }
    session = state["session"]
    if session:
        elapsed = elapsed_seconds(session, now)
        if session["kind"] == "assessment":
            assessment_seconds = session["assessment_minutes"] * 60
            plan = {"assessment": session["assessment_minutes"], "feedback": 15}
        else:
            plan = session.get("plan_minutes", study_plan(60))
        total_seconds = sum(plan.values()) * 60
        boundary = 0
        phase = "session time complete"
        for section, minutes in plan.items():
            boundary += minutes * 60
            if elapsed < boundary:
                phase = section
                break
        result["session"] = {
            "kind": session["kind"],
            "language": session.get("language", "Go"),
            "status": "paused" if is_paused(session) else "running",
            "owner_thread_id": session["thread_id"],
            "started_at": session["started_at"],
            "planned_minutes": sum(plan.values()),
            "plan_minutes": plan,
            "active_seconds": elapsed,
            "remaining_seconds": max(0, total_seconds - elapsed),
            "phase": phase,
            "phase_remaining_seconds": max(0, boundary - elapsed),
        }
        result["language"] = session.get("language", "Go")
        if session["kind"] == "assessment":
            result["session"]["assessment_remaining_seconds"] = max(0, assessment_seconds - elapsed)
            result.pop("next_task")
            result.pop("teaching_contract")
            result.pop("unit")
            result.pop("weekday_work")
    return "Course context at turn start:\n" + json.dumps(result, indent=2)


def change_session(state, args, now):
    session = state["session"]
    stamp = now.isoformat()
    if args.command == "start":
        if session:
            raise ValueError("A session is already open. Resume or end it before starting another.")
        if args.kind == "assessment" and args.minutes is not None:
            raise ValueError("Use --assessment-minutes for an assessment. --minutes is for ordinary study.")
        plan = study_plan(args.minutes) if args.kind == "study" else None
        state["session"] = {
            "kind": args.kind,
            "language": args.language,
            "assessment_minutes": args.assessment_minutes if args.kind == "assessment" else None,
            "thread_id": args.thread_id,
            "started_at": stamp,
            "ended_at": None,
            "pauses": [],
        }
        if plan:
            state["session"]["plan_minutes"] = plan
        return
    if not session:
        raise ValueError("No study session is open.")
    if args.command == "pause":
        if is_paused(session):
            raise ValueError("The session is already paused.")
        split_task = session["assessment_minutes"] == 90 and args.session_boundary
        if session["kind"] == "assessment" and not (args.infrastructure_failure or split_task):
            raise ValueError("Pause only for an infrastructure failure or a 90-minute task session boundary.")
        state["session"]["pauses"].append({"start": stamp, "end": None, "reason": args.reason})
    elif args.command == "resume":
        if not is_paused(session):
            raise ValueError("The session is already running.")
        session["pauses"][-1]["end"] = stamp
    elif args.command == "end":
        session["ended_at"] = stamp
        session["active_seconds"] = elapsed_seconds(session, now)
        if is_paused(session):
            session["pauses"][-1]["end"] = stamp
        session["outcome"] = "Review pending"
        state["history"].append(session)
        state["session"] = None


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("command", choices=["init", "context", "hook", "plan", "start", "pause", "resume", "end"])
    parser.add_argument("--state", type=Path, default=STATE_PATH)
    parser.add_argument("--kind", choices=["study", "assessment"], default="study")
    parser.add_argument("--language", choices=["Go", "C++", "Rust"], default="Go")
    parser.add_argument("--minutes", type=int, help="Available active minutes for an ordinary study session")
    parser.add_argument("--assessment-minutes", type=int, choices=[45, 60, 90], default=45)
    parser.add_argument("--thread-id", default="manual")
    parser.add_argument("--reason", default="User-reported break")
    parser.add_argument("--infrastructure-failure", action="store_true")
    parser.add_argument("--session-boundary", action="store_true")
    args = parser.parse_args()
    now = datetime.now(timezone.utc)
    try:
        if args.command == "plan":
            if args.kind == "assessment":
                raise ValueError("The plan command divides ordinary study time. Assessment limits stay fixed.")
            print(json.dumps(study_plan(args.minutes), indent=2))
            return
        state = read_state(args.state)
        if args.command == "hook":
            event = json.load(sys.stdin)
            output = context(state, now, args.state, event["session_id"])
            print(json.dumps({"hookSpecificOutput": {
                "hookEventName": event["hook_event_name"],
                "additionalContext": output,
            }}))
            return
        if args.command in ("start", "pause", "resume", "end"):
            change_session(state, args, now)
            save_state(args.state, state)
        elif args.command == "init" and not args.state.exists():
            save_state(args.state, state)
        print(context(state, now, args.state, args.thread_id))
    except (ValueError, KeyError, OSError) as error:
        parser.error(str(error))


if __name__ == "__main__":
    main()
