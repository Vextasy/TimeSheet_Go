package domain

import "time"

type TimesheetServices struct {
	// the interfaces used by Timesheet
	Graph GraphSvc    // Microsoft Graph client.
	Cal   CalendarSvc // Read tasks from Microsoft Graph Outlook calendar.
	Dump  DumpSvc     // Dump tasks to stdout.
}
type TimesheetSvc interface {
	Run() error
}

type GraphSvc interface {
	Read(userName string, fromDate time.Time, toDate time.Time) ([]Task, error)
}

type CalendarSvc interface {
	Aggregate(tasks []Task) []*Project
}

type DumpSvc interface {
	Projects(projects []*Project) []string
}

// Credentials required by the TokenAcquirerFactory.
//   - "Instance": "https://login.microsoftonline.com/",
//   - "TenantId": "<the azure tenant id>",
//   - "ClientId": "<the app registration application (client) id>",
//   - "ClientSecret": "<the azure client secret>"
type Auth struct {
	Instance     string
	TenantId     string
	ClientId     string
	ClientSecret string
}
