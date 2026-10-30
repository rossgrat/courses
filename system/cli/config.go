package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/url"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
)

type source struct {
	Path   string `json:"path,omitempty"`
	URL    string `json:"url,omitempty"`
	Format string `json:"format"`
}

type learningPolicy struct {
	Driver         string   `json:"driver"`
	AgentRole      string   `json:"agent_role"`
	Explanations   string   `json:"explanations"`
	Implementation string   `json:"implementation"`
	Guidance       []string `json:"guidance,omitempty"`
}

type sessionPolicy struct {
	Timing string `json:"timing"`
	Split  struct {
		Retrieval   int `json:"retrieval"`
		FocusedWork int `json:"focused_work"`
		Feedback    int `json:"feedback"`
	} `json:"split,omitempty"`
}

type reviewPolicy struct {
	IntervalDays []int `json:"interval_days"`
	WeekdaysOnly bool  `json:"weekdays_only"`
}

type assessmentPolicy struct {
	Minutes      []int `json:"minutes"`
	Feedback     int   `json:"feedback"`
	SplitMinutes []int `json:"split_minutes"`
}

type course struct {
	SchemaVersion int              `json:"schema_version"`
	ID            string           `json:"id"`
	Name          string           `json:"name"`
	Sources       []source         `json:"sources"`
	Learning      learningPolicy   `json:"learning"`
	Session       *sessionPolicy   `json:"session,omitempty"`
	Reviews       *reviewPolicy    `json:"reviews,omitempty"`
	Assessment    assessmentPolicy `json:"assessment,omitempty"`
	Initial       object           `json:"initial,omitempty"`
	WeekdayWork   []string         `json:"weekday_work,omitempty"`

	dir, root string
}

func repositoryRoot(cwd string) (string, error) {
	output, err := exec.Command("git", "-C", cwd, "rev-parse", "--show-toplevel").Output()
	if err != nil {
		return "", errors.New("the course command requires a Git working tree")
	}
	return strings.TrimSpace(string(output)), nil
}

func findCourse(cwd, selected string) (*course, error) {
	root, err := repositoryRoot(cwd)
	if err != nil {
		return nil, err
	}
	if selected != "" {
		if filepath.Base(selected) != selected || selected == "." || selected == ".." {
			return nil, errors.New("course must be a course identifier")
		}
		return loadCourse(filepath.Join(root, "programs", selected), root)
	}
	for dir := cwd; ; dir = filepath.Dir(dir) {
		if filepath.Dir(dir) == filepath.Join(root, "programs") {
			return loadCourse(dir, root)
		}
		if dir == root || filepath.Dir(dir) == dir {
			return nil, nil
		}
	}
}

func loadCourse(dir, root string) (*course, error) {
	path := filepath.Join(dir, "course.json")
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	c := course{dir: dir, root: root}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	decoder.UseNumber()
	if err := decoder.Decode(&c); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	if err := decoder.Decode(new(any)); err != io.EOF {
		return nil, fmt.Errorf("%s must contain one JSON object", path)
	}
	if err := c.validate(); err != nil {
		return nil, fmt.Errorf("%s: %w", path, err)
	}
	return &c, nil
}

func localFile(root, path string) error {
	if !filepath.IsLocal(path) || strings.ContainsAny(path, "\r\n") {
		return fmt.Errorf("%q must be a relative file path within its folder", path)
	}
	info, err := os.Stat(filepath.Join(root, path))
	if err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return fmt.Errorf("%q must be a file", path)
	}
	return nil
}

func (c *course) validate() error {
	if c.SchemaVersion != 1 {
		return errors.New("schema_version must be 1")
	}
	if c.ID != filepath.Base(c.dir) || strings.TrimSpace(c.Name) == "" {
		return errors.New("id must match the course folder and name must be nonempty")
	}
	if len(c.Sources) == 0 {
		return errors.New("sources must contain at least one document or URL")
	}
	for _, s := range c.Sources {
		if (s.Path == "") == (s.URL == "") {
			return errors.New("each source requires exactly one of path or url")
		}
		if !slices.Contains([]string{"latex", "markdown", "pdf", "text", "web"}, s.Format) {
			return fmt.Errorf("unknown source format %q", s.Format)
		}
		if s.Path != "" {
			if err := localFile(c.dir, s.Path); err != nil {
				return err
			}
			if s.Format == "web" || s.Format == "latex" && filepath.Ext(s.Path) != ".tex" {
				return errors.New("web sources require a URL and latex sources require a .tex file")
			}
		} else {
			u, err := url.ParseRequestURI(s.URL)
			if err != nil || u.Host == "" || !slices.Contains([]string{"http", "https"}, u.Scheme) || s.Format == "latex" {
				return errors.New("source URLs require HTTP or HTTPS; latex sources must be local")
			}
		}
	}
	if !slices.Contains([]string{"curriculum", "book", "project", "learner"}, c.Learning.Driver) {
		return errors.New("learning.driver must be curriculum, book, project, or learner")
	}
	if !slices.Contains([]string{"tutor", "collaborator", "consultant"}, c.Learning.AgentRole) {
		return errors.New("learning.agent_role must be tutor, collaborator, or consultant")
	}
	for _, policy := range []string{c.Learning.Explanations, c.Learning.Implementation} {
		if !slices.Contains([]string{"on_request", "proactive", "never"}, policy) {
			return errors.New("explanations and implementation must be on_request, proactive, or never")
		}
	}
	for _, path := range c.Learning.Guidance {
		if err := localFile(c.root, path); err != nil {
			return err
		}
	}
	if c.Session != nil {
		s := c.Session
		if !slices.Contains([]string{"budgeted", "untimed"}, s.Timing) {
			return errors.New("session.timing must be budgeted or untimed")
		}
		if s.Split.Retrieval < 0 || s.Split.Feedback < 0 || s.Split.FocusedWork < 0 {
			return errors.New("session split weights cannot be negative")
		}
		if s.Timing == "budgeted" && s.Split.FocusedWork == 0 {
			return errors.New("budgeted sessions require a positive focused_work weight")
		}
		if s.Timing == "untimed" && s.Split.Retrieval+s.Split.Feedback+s.Split.FocusedWork != 0 {
			return errors.New("untimed sessions cannot define a split")
		}
	}
	if c.Reviews != nil && !increasingPositive(c.Reviews.IntervalDays) {
		return errors.New("review interval_days must be positive and strictly increasing")
	}
	a := c.Assessment
	if len(a.Minutes) != 0 || a.Feedback != 0 || len(a.SplitMinutes) != 0 {
		if !increasingPositive(a.Minutes) || a.Feedback < 1 {
			return errors.New("assessment minutes must be positive and increasing, with at least one feedback minute")
		}
		for _, duration := range a.SplitMinutes {
			if !slices.Contains(a.Minutes, duration) {
				return errors.New("assessment split_minutes must be allowed assessment durations")
			}
		}
	}
	if len(c.WeekdayWork) != 0 && len(c.WeekdayWork) != 7 {
		return errors.New("weekday_work must contain seven entries, starting with Monday")
	}
	for key, value := range c.Initial {
		if key == "stage" {
			n, ok := value.(json.Number)
			stage, err := n.Int64()
			if !ok || err != nil || stage < 0 {
				return errors.New("initial.stage must be a nonnegative integer")
			}
		} else {
			if !slices.Contains([]string{"language", "unit", "position", "next_task"}, key) {
				return fmt.Errorf("unknown initial field %q", key)
			}
			if _, ok := value.(string); !ok {
				return fmt.Errorf("initial.%s must be text", key)
			}
		}
	}
	return nil
}

func increasingPositive(values []int) bool {
	previous := 0
	for _, value := range values {
		if value <= previous {
			return false
		}
		previous = value
	}
	return len(values) != 0
}

func courseCatalog(cwd, selected, command string) error {
	root, err := repositoryRoot(cwd)
	if err != nil {
		return err
	}
	paths, err := filepath.Glob(filepath.Join(root, "programs", "*", "course.json"))
	if err != nil {
		return err
	}
	if selected != "" {
		c, err := findCourse(cwd, selected)
		if err != nil {
			return err
		}
		paths = []string{filepath.Join(c.dir, "course.json")}
	}
	entries := []object{}
	pdfs := []string{}
	for _, path := range paths {
		c, err := loadCourse(filepath.Dir(path), root)
		if err != nil {
			return err
		}
		entries = append(entries, object{"id": c.ID, "name": c.Name, "driver": c.Learning.Driver, "directory": c.dir})
		for _, s := range c.Sources {
			if s.Format == "latex" {
				relative, err := filepath.Rel(root, filepath.Join(c.dir, s.Path))
				if err != nil {
					return err
				}
				pdfs = append(pdfs, relative)
			}
		}
	}
	if command == "pdf-sources" {
		for _, path := range pdfs {
			fmt.Println(path)
		}
		return nil
	}
	return printJSON(entries)
}
