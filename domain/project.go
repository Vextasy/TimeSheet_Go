package domain

import (
	"slices"
	"time"
)

type Project struct {
	Name    string
	Tasks   []Task
	Summary map[string]*TaskSummary
	Groups  map[string]*GroupSummary
}

// A Task represents an Outlook event that has
// the format "Project - Group - Description".
type Task struct {
	Project  string
	Group    string
	Desc     string
	Start    time.Time
	Duration time.Duration
}

// Within a Project a TaskSummary is a summary of all tasks
// that have the same group and description.
type TaskSummary struct {
	Group    string
	Desc     string
	Duration time.Duration
	Started  time.Time // earliest task start
}

// Within a Project a GroupSummary is a summary of all tasks
// that have the same group.
type GroupSummary struct {
	Group    string
	Duration time.Duration
	Started  time.Time // earliest task start
}

func NewProject(name string) *Project {
	return &Project{
		Name:    name,
		Tasks:   []Task{},
		Summary: map[string]*TaskSummary{},
		Groups:  map[string]*GroupSummary{},
	}
}

func (proj *Project) AddTask(task Task) {
	proj.Tasks = append(proj.Tasks, task)
}

func newTaskSummary(desc string, group string) *TaskSummary {
	return &TaskSummary{
		Group:    group,
		Desc:     desc,
		Duration: 0,
	}
}

func newGroupSummary(group string) *GroupSummary {
	return &GroupSummary{
		Group:    group,
		Duration: 0,
	}
}

// Add the task to the appropriate project summary and groups
// and increment their durations.
func (proj *Project) sum(t Task) {
	// Project Summary is desc within group.
	key := t.Group + "-" + t.Desc
	if _, ok := proj.Summary[key]; !ok {
		proj.Summary[key] = newTaskSummary(t.Desc, t.Group)
	}
	s := proj.Summary[key]
	s.Duration += t.Duration
	if s.Started.IsZero() || s.Started.After(t.Start) {
		s.Started = t.Start
	}

	// Project Group is just group.
	key = t.Group
	if _, ok := proj.Groups[key]; !ok {
		proj.Groups[key] = newGroupSummary(t.Group)
	}
	g := proj.Groups[key]
	g.Duration += t.Duration
	if g.Started.IsZero() || g.Started.After(t.Start) {
		g.Started = t.Start
	}
}

// Arrange tasks in start order and summarize.
func (proj *Project) Summarize() {
	slices.SortFunc(proj.Tasks, func(a, b Task) int {
		return a.Start.Compare(b.Start)
	})
	for _, task := range proj.Tasks {
		proj.sum(task)
	}
}
