package svc

import (
	"github.com/vextasy/Timesheet_go/domain"
)

func NewServices(cfg TsConfig) domain.TimesheetServices {
	return domain.TimesheetServices{
		Graph: NewGraphSvc(cfg.Auth),
		Cal:   NewCalendarSvc(),
		Dump:  NewDumpSvc(cfg),
	}
}
