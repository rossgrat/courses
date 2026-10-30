package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"time"
	_ "time/tzdata"
)

type object = map[string]any

type state struct {
	Progress object   `json:"progress"`
	Reviews  []object `json:"reviews"`
	Session  object   `json:"session"`
	History  []object `json:"history"`
}

type options struct {
	command, selected, statePath, kind, language, thread, reason string
	minutes, assessment                                          int
	infrastructure, boundary                                     bool
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "course:", err)
		os.Exit(1)
	}
}

func run() error {
	if len(os.Args) < 2 {
		return errors.New("use list, validate, pdf-sources, context, hook, plan, init, start, pause, resume, or end")
	}
	o := options{command: os.Args[1]}
	flags := flag.NewFlagSet("course "+o.command, flag.ContinueOnError)
	flags.StringVar(&o.selected, "course", "", "Course identifier when outside its folder")
	flags.StringVar(&o.statePath, "state", "", "Alternate state file for a manual walkthrough")
	flags.StringVar(&o.kind, "kind", "study", "study or assessment")
	flags.StringVar(&o.language, "language", "", "Language for this session")
	flags.StringVar(&o.thread, "thread-id", "manual", "Current conversation identifier")
	flags.StringVar(&o.reason, "reason", "User-reported break", "Reason for a pause")
	flags.IntVar(&o.minutes, "minutes", 0, "Available active study minutes")
	flags.IntVar(&o.assessment, "assessment-minutes", 0, "Course assessment limit")
	flags.BoolVar(&o.infrastructure, "infrastructure-failure", false, "Pause for an infrastructure failure")
	flags.BoolVar(&o.boundary, "session-boundary", false, "Pause an assessment that permits multiple sessions")
	if err := flags.Parse(os.Args[2:]); err != nil {
		return err
	}
	if flags.NArg() != 0 {
		return errors.New("unexpected positional argument")
	}
	if !slices.Contains([]string{"list", "validate", "pdf-sources", "context", "hook", "plan", "init", "start", "pause", "resume", "end"}, o.command) {
		return fmt.Errorf("unknown command %q", o.command)
	}
	if o.kind != "study" && o.kind != "assessment" {
		return errors.New("kind must be study or assessment")
	}
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}
	if slices.Contains([]string{"list", "validate", "pdf-sources"}, o.command) {
		return courseCatalog(cwd, o.selected, o.command)
	}
	var event struct {
		Name    string `json:"hook_event_name"`
		Session string `json:"session_id"`
		CWD     string `json:"cwd"`
	}
	if o.command == "hook" {
		if err := json.NewDecoder(os.Stdin).Decode(&event); err != nil {
			return err
		}
		if event.CWD != "" {
			cwd = event.CWD
		}
		o.thread = event.Session
	}
	c, err := findCourse(cwd, o.selected)
	if err != nil {
		return err
	}
	if c == nil {
		if o.command == "hook" {
			return nil
		}
		if o.command == "context" {
			return printJSON(object{"course": nil, "message": "Open a course folder or pass --course COURSE."})
		}
		return errors.New("open a course folder or pass --course COURSE")
	}
	if o.command == "plan" {
		if o.kind != "study" {
			return errors.New("plan divides ordinary study time. Assessment limits stay fixed")
		}
		plan, err := studyPlan(c, o.minutes)
		if err != nil {
			return err
		}
		return printJSON(plan)
	}
	if o.statePath == "" {
		o.statePath = filepath.Join(c.dir, "state.json")
	}
	o.statePath, err = filepath.Abs(o.statePath)
	if err != nil {
		return err
	}
	s, exists, err := readState(o.statePath)
	if err != nil {
		return err
	}
	if !exists && (o.command == "init" || o.command == "start") {
		for key, value := range c.Initial {
			s.Progress[key] = value
		}
	}
	now := time.Now().UTC()
	switch o.command {
	case "start", "pause", "resume", "end":
		if err := changeSession(&s, c, o, now); err != nil {
			return err
		}
		if err := saveJSON(o.statePath, s); err != nil {
			return err
		}
		exists = true
	case "init":
		if !exists {
			if err := saveJSON(o.statePath, s); err != nil {
				return err
			}
			exists = true
		}
	}
	context, err := courseContext(s, c, o, now, exists)
	if err != nil {
		return err
	}
	if o.command == "hook" {
		data, err := json.MarshalIndent(context, "", "  ")
		if err != nil {
			return err
		}
		return printJSON(object{"hookSpecificOutput": object{
			"hookEventName":     event.Name,
			"additionalContext": "Course context at turn start:\n" + string(data),
		}})
	}
	return printJSON(context)
}

func readState(path string) (state, bool, error) {
	data, err := os.ReadFile(path)
	if errors.Is(err, os.ErrNotExist) {
		return state{
			Progress: object{"evidence": []any{}},
			Reviews:  []object{}, History: []object{},
		}, false, nil
	}
	if err != nil {
		return state{}, false, err
	}
	var s state
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(&s); err != nil {
		return s, true, fmt.Errorf("%s: %w", path, err)
	}
	if s.Progress == nil || s.Reviews == nil || s.History == nil {
		return s, true, errors.New("state requires progress, reviews, and history. The existing file was not changed")
	}
	return s, true, nil
}

func saveJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return err
	}
	return writeAtomic(path, append(data, '\n'))
}

func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	file, err := os.CreateTemp(filepath.Dir(path), ".course-*")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()
	if _, err := file.Write(data); err != nil {
		return err
	}
	if err := file.Sync(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), path)
}

func printJSON(value any) error {
	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "  ")
	return encoder.Encode(value)
}

func text(m object, key string) string {
	s, _ := m[key].(string)
	return s
}

func number(m object, key string) int {
	switch n := m[key].(type) {
	case int:
		return n
	case json.Number:
		value, _ := n.Int64()
		return int(value)
	}
	return 0
}
