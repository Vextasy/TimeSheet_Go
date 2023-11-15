package svc

import (
	"fmt"
	"strings"
	"time"

	"github.com/vextasy/Timesheet_go/domain"
)

type TsConfig struct {
	UserName string
	DateFrom time.Time
	DateTo   time.Time
	Auth     domain.Auth
}

// tsSvc implements domain.TimesheetSvc.
type tsSvc struct {
	domain.TimesheetServices
	cfg TsConfig
}

func NewTimesheetSvc(cfg TsConfig, svc domain.TimesheetServices) domain.TimesheetSvc {
	return tsSvc{
		cfg:               cfg,
		TimesheetServices: svc,
	}
}

func (svc tsSvc) Run() error {
	tasks, err := svc.Graph.Read(svc.cfg.UserName, svc.cfg.DateFrom, svc.cfg.DateTo)
	if err != nil {
		return err
	}
	projects := svc.Cal.Aggregate(tasks)
	lines := svc.Dump.Projects(projects)
	fmt.Println(strings.Join(lines, "\n"))
	return nil
}
