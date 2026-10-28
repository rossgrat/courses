package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"slices"
	"time"
)

func studyPlan(minutes int) (object, error) {
	if minutes < 3 {
		return nil, errors.New("give at least three active study minutes with --minutes")
	}
	review := max(1, (minutes+3)/6)
	return object{"retrieval": review, "focused work": minutes - 2*review, "feedback": review}, nil
}

func paused(session object) bool {
	pauses, _ := session["pauses"].([]any)
	if len(pauses) == 0 {
		return false
	}
	last, _ := pauses[len(pauses)-1].(object)
	return last != nil && last["end"] == nil
}

func elapsed(session object, now time.Time) (int, error) {
	start, err := time.Parse(time.RFC3339Nano, text(session, "started_at"))
	if err != nil {
		return 0, err
	}
	end := now
	if stamp := text(session, "ended_at"); stamp != "" {
		end, err = time.Parse(time.RFC3339Nano, stamp)
		if err != nil {
			return 0, err
		}
	}
	duration := end.Sub(start)
	pauses, _ := session["pauses"].([]any)
	for _, item := range pauses {
		pause, ok := item.(object)
		if !ok {
			return 0, errors.New("invalid pause record")
		}
		begin, err := time.Parse(time.RFC3339Nano, text(pause, "start"))
		if err != nil {
			return 0, err
		}
		finish := end
		if stamp := text(pause, "end"); stamp != "" {
			finish, err = time.Parse(time.RFC3339Nano, stamp)
			if err != nil {
				return 0, err
			}
		}
		duration -= finish.Sub(begin)
	}
	if duration < 0 {
		return 0, errors.New("session timestamps are in the future. Make sure that the machine clock is correct")
	}
	return int(duration.Seconds()), nil
}

func changeSession(s *state, c *course, o options, now time.Time) error {
	stamp := now.Format(time.RFC3339Nano)
	if o.command == "start" {
		if s.Session != nil {
			return errors.New("a session is already open. Resume or end it before starting another")
		}
		language := o.language
		if language == "" {
			language = c.Language
		}
		var plan object
		var assessment any
		if o.kind == "assessment" {
			if o.minutes != 0 {
				return errors.New("use --assessment-minutes for an assessment. --minutes is for ordinary study")
			}
			if len(c.Assessment.Minutes) == 0 {
				return errors.New("this course does not define timed assessments")
			}
			if o.assessment == 0 {
				o.assessment = c.Assessment.Minutes[0]
			}
			if !slices.Contains(c.Assessment.Minutes, o.assessment) {
				return fmt.Errorf("assessment minutes must be one of %v", c.Assessment.Minutes)
			}
			assessment = o.assessment
			plan = object{"assessment": o.assessment, "feedback": c.Assessment.Feedback}
		} else {
			if o.assessment != 0 {
				return errors.New("--assessment-minutes requires --kind assessment")
			}
			var err error
			plan, err = studyPlan(o.minutes)
			if err != nil {
				return err
			}
		}
		s.Session = object{
			"kind": o.kind, "language": language, "assessment_minutes": assessment,
			"thread_id": o.thread, "started_at": stamp, "ended_at": nil, "pauses": []any{}, "plan_minutes": plan,
		}
		return nil
	}
	if s.Session == nil {
		return errors.New("no study session is open")
	}
	active, err := elapsed(s.Session, now)
	if err != nil {
		return err
	}
	pauses, _ := s.Session["pauses"].([]any)
	switch o.command {
	case "pause":
		if paused(s.Session) {
			return errors.New("the session is already paused")
		}
		split := o.boundary && slices.Contains(c.Assessment.SplitMinutes, number(s.Session, "assessment_minutes"))
		if text(s.Session, "kind") == "assessment" && !o.infrastructure && !split {
			return errors.New("pause this assessment only for an infrastructure failure or an allowed session boundary")
		}
		s.Session["pauses"] = append(pauses, object{"start": stamp, "end": nil, "reason": o.reason})
	case "resume":
		if !paused(s.Session) {
			return errors.New("the session is already running")
		}
		pauses[len(pauses)-1].(object)["end"] = stamp
		s.Session["thread_id"] = o.thread
	case "end":
		if paused(s.Session) {
			pauses[len(pauses)-1].(object)["end"] = stamp
		}
		s.Session["ended_at"] = stamp
		s.Session["active_seconds"] = active
		s.Session["outcome"] = "Review pending"
		s.History = append(s.History, s.Session)
		s.Session = nil
	}
	return nil
}

func courseContext(s state, c *course, o options, now time.Time, exists bool) (object, error) {
	zone, err := time.LoadLocation("America/Chicago")
	if err != nil {
		return nil, err
	}
	local := now.In(zone)
	due := []string{}
	for _, review := range s.Reviews {
		if text(review, "status") != "retired" && text(review, "due") <= local.Format("2006-01-02") {
			due = append(due, text(review, "id"))
		}
	}
	result := object{
		"course": c.Name, "course_id": c.ID, "local_time": local.Format(time.RFC3339), "timezone": zone.String(),
		"weekday": local.Weekday().String(), "thread_id": o.thread, "state_file": o.statePath, "state_exists": exists,
		"source": filepath.Join(c.dir, c.Source), "stage": s.Progress["stage"], "unit": s.Progress["unit"],
		"language": s.Progress["language"], "next_task": s.Progress["next_task"],
		"reviews_due": len(due), "due_review_ids": due[:min(len(due), 10)], "session": "inactive",
	}
	if contract := text(s.Progress, "teaching_contract"); contract != "" {
		result["teaching_contract"] = contract
	}
	if len(c.WeekdayWork) == 7 {
		result["weekday_work"] = c.WeekdayWork[(int(local.Weekday())+6)%7]
	}
	if s.Session == nil {
		return result, nil
	}
	active, err := elapsed(s.Session, now)
	if err != nil {
		return nil, err
	}
	plan, _ := s.Session["plan_minutes"].(object)
	order := []string{"retrieval", "focused work", "feedback"}
	assessment := text(s.Session, "kind") == "assessment"
	if assessment {
		order = []string{"assessment", "feedback"}
		if plan == nil {
			plan = object{"assessment": number(s.Session, "assessment_minutes"), "feedback": c.Assessment.Feedback}
		}
	} else if plan == nil {
		plan, _ = studyPlan(60)
	}
	total, boundary, remaining := 0, 0, 0
	phase := "session time complete"
	for _, section := range order {
		total += number(plan, section)
		boundary += number(plan, section) * 60
		if phase == "session time complete" && active < boundary {
			phase, remaining = section, boundary-active
		}
	}
	status := "running"
	if paused(s.Session) {
		status = "paused"
	}
	summary := object{
		"kind": s.Session["kind"], "language": s.Session["language"], "status": status,
		"owner_thread_id": s.Session["thread_id"], "started_at": s.Session["started_at"],
		"planned_minutes": total, "plan_minutes": plan, "active_seconds": active,
		"remaining_seconds": max(0, total*60-active), "phase": phase, "phase_remaining_seconds": remaining,
	}
	result["session"] = summary
	result["language"] = s.Session["language"]
	if assessment {
		summary["assessment_remaining_seconds"] = max(0, number(s.Session, "assessment_minutes")*60-active)
		for _, key := range []string{"next_task", "unit", "weekday_work", "due_review_ids", "teaching_contract"} {
			delete(result, key)
		}
	}
	return result, nil
}
