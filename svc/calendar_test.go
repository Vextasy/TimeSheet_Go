package svc

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vextasy/Timesheet_go/domain"
)

var svc domain.CalendarSvc

const hr time.Duration = time.Hour
const min time.Duration = time.Minute

func init() {
	svc = NewCalendarSvc()
}

// Creates a map of projects with their tasks and summary.
// Projects are sorted by their first task start time.
func Test_creates_map_of_projects(t *testing.T) {
	tasks := []domain.Task{
		{Project: "Project 1", Start: time.Now(), Duration: hr},
		{Project: "Project 2", Start: time.Now(), Duration: hr},
		{Project: "Project 1", Start: time.Now(), Duration: hr},
	}
	projects := svc.Aggregate(tasks)
	assert.Equal(t, 2, len(projects))
	assert.NotNil(t, projects[0].Summary)
	assert.NotNil(t, projects[1].Summary)
	assert.Equal(t, "Project 1", projects[0].Name)
	assert.Equal(t, "Project 2", projects[1].Name)
}

// Adds tasks to existing projects.
func Test_adds_tasks_to_existing_projects(t *testing.T) {
	tasks := []domain.Task{
		{Project: "Project 1", Start: time.Now(), Duration: hr},
		{Project: "Project 2", Start: time.Now(), Duration: hr},
		{Project: "Project 1", Start: time.Now(), Duration: hr},
	}
	projects := svc.Aggregate(tasks)
	assert.Equal(t, 2, len(projects))
	assert.Equal(t, 2, len(projects[0].Tasks))
	assert.Equal(t, 1, len(projects[1].Tasks))
}

// Summarizes each project after adding tasks.
func Test_summarizes_each_project(t *testing.T) {
	tasks := []domain.Task{
		{Project: "Project 1", Group: "Group 1", Desc: "Desc 1", Start: time.Now(), Duration: hr},
		{Project: "Project 2", Group: "Group 1", Desc: "Desc 1", Start: time.Now(), Duration: hr},
		{Project: "Project 1", Group: "Group 1", Desc: "Desc 2", Start: time.Now(), Duration: 15 * min},
	}
	projects := svc.Aggregate(tasks)
	assert.Equal(t, 2, len(projects))
	assert.NotNil(t, projects[0].Summary)
	assert.NotNil(t, projects[1].Summary)
	assert.NotNil(t, projects[0].Groups)
	assert.NotNil(t, projects[1].Groups)
	assert.Equal(t, 1*hr+15*min, projects[0].Groups["Group 1"].Duration)
	assert.Equal(t, 1*hr, projects[0].Summary["Group 1-Desc 1"].Duration)
	assert.Equal(t, 15*min, projects[0].Summary["Group 1-Desc 2"].Duration)
}

// Handles empty input task list.
func Test_handles_empty_input_task_list(t *testing.T) {
	tasks := []domain.Task{}
	projects := svc.Aggregate(tasks)
	assert.Equal(t, 0, len(projects))
}

// Handles tasks with empty project name.
func Test_handles_tasks_with_empty_project_name(t *testing.T) {
	tasks := []domain.Task{
		{Project: "", Start: time.Now(), Duration: hr},
		{Project: "Project 2", Start: time.Now(), Duration: hr},
		{Project: "", Start: time.Now(), Duration: hr},
	}
	projects := svc.Aggregate(tasks)
	assert.Equal(t, 2, len(projects))
	assert.Equal(t, 2, len(projects[0].Tasks))
	assert.Equal(t, 1, len(projects[1].Tasks))
	assert.Equal(t, "", projects[0].Name)
}

func Test_calendarSvc_Aggregate(t *testing.T) {
	start := time.Now()
	type args struct {
		tasks []domain.Task
	}
	tests := []struct {
		name string
		args args
		want []*domain.Project
	}{
		{
			"Test 1",
			args{tasks: []domain.Task{
				{Project: "Project 1", Group: "Group 1", Desc: "Desc 1", Start: start, Duration: hr},
				{Project: "Project 1", Group: "Group 1", Desc: "Desc 2", Start: start, Duration: 30 * min},
			}},
			[]*domain.Project{
				{
					Name: "Project 1",
					Tasks: []domain.Task{
						{Project: "Project 1", Group: "Group 1", Desc: "Desc 1", Start: start, Duration: hr},
						{Project: "Project 1", Group: "Group 1", Desc: "Desc 2", Start: start, Duration: 30 * min},
					},
					Summary: map[string]*domain.TaskSummary{
						"Group 1-Desc 1": {
							Group:    "Group 1",
							Desc:     "Desc 1",
							Duration: hr,
							Started:  start,
						},
						"Group 1-Desc 2": {
							Group:    "Group 1",
							Desc:     "Desc 2",
							Duration: 30 * min,
							Started:  start,
						},
					},
					Groups: map[string]*domain.GroupSummary{
						"Group 1": {
							Group:    "Group 1",
							Duration: hr + 30*min,
							Started:  start,
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			svc := calendarSvc{}
			if got := svc.Aggregate(tt.args.tasks); !reflect.DeepEqual(got, tt.want) {
				jgot, _ := json.Marshal(got)
				jwant, _ := json.Marshal(tt.want)
				t.Errorf("\ngot:  %s\nwant: %s", jgot, jwant)
			}
		})
	}
}
