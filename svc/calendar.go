package svc

import (
	"slices"

	"github.com/vextasy/Timesheet_go/domain"
)

// calendarSvc implements domain.CalendarSvc.
type calendarSvc struct {
}

func NewCalendarSvc() domain.CalendarSvc {
	return calendarSvc{}
}

func (svc calendarSvc) Aggregate(tasks []domain.Task) []*domain.Project {
	p := make(map[string]*domain.Project) // Map by task.Project string.
	for _, task := range tasks {
		if _, ok := p[task.Project]; !ok {
			p[task.Project] = domain.NewProject(task.Project)
		}
		p[task.Project].AddTask(task)
	}
	projects := make([]*domain.Project, 0, len(p))
	for _, project := range p {
		project.Summarize()
		projects = append(projects, project)
	}
	// Arrange the projects in time order of their first task.
	slices.SortFunc(projects, func(a, b *domain.Project) int { return a.Tasks[0].Start.Compare(b.Tasks[0].Start) })
	return projects
}
