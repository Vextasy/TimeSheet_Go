package svc

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/vextasy/Timesheet_go/domain"
)

// dumpSvc implements domain.DumpSvc.
type dumpSvc struct {
	cfg TsConfig
}

func NewDumpSvc(cfg TsConfig) domain.DumpSvc {
	return dumpSvc{cfg}
}

func (svc dumpSvc) Projects(projects []*domain.Project) []string {
	output := []string{}

	// Return the sum of the durations of the tasks in the task list.
	Time := func(tl []domain.Task) time.Duration {
		duration := time.Duration(0)
		for _, t := range tl {
			duration += t.Duration
		}
		return duration
	}

	// Return a string that represents the date part of t.
	// Such strings can be compared for equality to determine
	// if two time.Time values represent the same date.
	DateStr := func(t time.Time) string {
		return t.Format("2006-01-02")
	}

	// Return the sum of the durations of the tasks in the task list
	// that started on the given day.
	DayTime := func(tl []domain.Task, day time.Time) time.Duration {
		duration := time.Duration(0)
		for _, t := range tl {
			if DateStr(t.Start) == DateStr(day) {
				duration += t.Duration
			}
		}
		return duration
	}

	dmyFormat := "02 Jan 2006"
	output = append(output, fmt.Sprintf("For the Dates %s - %s", svc.cfg.DateFrom.Format(dmyFormat), svc.cfg.DateTo.Format(dmyFormat)))
	output = append(output, "")

	for _, proj := range projects {
		output = append(output, fmt.Sprintf("%s = %s", proj.Name, fmtLongTime(Time(proj.Tasks))))

		startdays := mondays(svc.cfg.DateFrom)
		for _, d := range startdays {
			var times []time.Duration
			var stimes []string // formatted string times.
			for _, day := range weekdays(d) {
				dt := DayTime(proj.Tasks, day)
				times = append(times, dt)
				stimes = append(stimes, fmtTime(dt))
			}
			eqn := fmt.Sprintf("%s", strings.Join(stimes, " + "))

			var sum time.Duration
			for _, t := range times {
				sum += t
			}
			output = append(output, fmt.Sprintf("w/b %s - %s = %s", fmtDate(d), eqn, fmtTime(sum)))
		}

		// Output the groups in order of start time of the earliest task.
		groups := make([]*domain.GroupSummary, 0, len(proj.Groups))
		for g := range proj.Groups {
			groups = append(groups, proj.Groups[g])
		}
		slices.SortFunc(groups, func(a, b *domain.GroupSummary) int { return a.Started.Compare(b.Started) })
		output = append(output, "")
		for _, g := range groups {
			output = append(output, fmt.Sprintf("- %s (%s)", g.Group, fmtLongTime(g.Duration)))
		}

		// Output the task summaries in order of the start time of the earliest task.
		tasks := make([]*domain.TaskSummary, 0, len(proj.Summary))
		for s := range proj.Summary {
			tasks = append(tasks, proj.Summary[s])
		}
		slices.SortFunc(tasks, func(a, b *domain.TaskSummary) int { return a.Started.Compare(b.Started) })
		output = append(output, "")
		for _, s := range tasks {
			output = append(output, fmt.Sprintf("- %s %s (%s)", s.Group, s.Desc, fmtLongTime(s.Duration)))
		}
		output = append(output, "")
	}

	return output
}

func fmtDate(d time.Time) string {
	day := d.Day()
	mon := d.Month()
	yr := d.Year()
	thisyr := time.Now().Year()

	if yr == thisyr {
		return fmt.Sprintf("%02d/%02d", day, mon)
	} else {
		return fmt.Sprintf("%02d/%02d/%d", day, mon, yr)
	}
}

// Return the number of hours and (remaining) minutes within a duration.
func hoursAndMinutes(ts time.Duration) (hours int, minutes int) {
	m := int(ts.Minutes())
	hours = m / 60
	minutes = m % 60
	return
}

func fmtLongTime(ts time.Duration) string {
	hrs, min := hoursAndMinutes(ts)

	if hrs == 0 && min == 0 {
		return "0"
	} else if hrs == 0 {
		return fmt.Sprintf("%d min", min)
	} else if min == 0 {
		return fmt.Sprintf("%d hr", hrs)
	} else {
		return fmt.Sprintf("%d hr %d min", hrs, min)
	}
}

func fmtTime(ts time.Duration) string {
	hrs, min := hoursAndMinutes(ts)

	if hrs == 0 && min == 0 {
		return "0"
	} else if min == 0 {
		return fmt.Sprintf("%d", hrs)
	} else {
		return fmt.Sprintf("%d:%d", hrs, min)
	}
}

// Return the sequence of mondays that contain days from the same month as `from'.
func mondays(from time.Time) []time.Time {
	dow := from.Weekday() // Sunday = 0
	dec := (int(dow) + 6) % 7
	mon1 := from.AddDate(0, 0, -dec)

	var mondays []time.Time
	for w := 0; w <= 5; w++ {
		mon := mon1.AddDate(0, 0, w*7)
		if mon.Month() != from.Month()+1 {
			mondays = append(mondays, mon)
		}
	}

	return mondays
}

// Return the sequence of saturdays that contain days from the same month as `from'
func saturdays(from time.Time) []time.Time {
	dow := from.Weekday()
	dec := int((dow + 1) % 7)
	sat1 := from.AddDate(0, 0, -dec)

	var result []time.Time

	for w := 0; w <= 5; w++ {
		sat := sat1.AddDate(0, 0, w*7)
		if sat.Month() != from.Month()+1 {
			result = append(result, sat)
		}
	}

	return result
}

func weekdays(monday time.Time) []time.Time {
	var result []time.Time

	for d := 0; d <= 6; d++ {
		result = append(result, monday.AddDate(0, 0, d))
	}

	return result
}
