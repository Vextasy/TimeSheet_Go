package main

import (
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/vextasy/Timesheet_go/internal/envrc"
	"github.com/vextasy/Timesheet_go/svc"
)

func main() {
	env := envrc.NewEnvRc(".")
	cfg := svc.TsConfig{}

	cfg.Auth.Instance = env.Get("Instance")
	cfg.Auth.TenantId = env.Get("TenantId")
	cfg.Auth.ClientId = env.Get("ClientId")
	cfg.Auth.ClientSecret = env.Get("ClientSecret")

	var nFlag = flag.Int("n", 0, "Produce a time sheet for 'n' months back.")
	var fromFlag = flag.String("from", "", "Override 'from' date (inclusive).")
	var toFlag = flag.String("to", "", "Override 'to' date (inclusive).")
	flag.StringVar(&cfg.UserName, "user", env.Get("UserName"), "User name.")
	flag.Parse()

	// Determine the date range from the -n flag (or its default).
	cfg.DateFrom, cfg.DateTo = monthOffset(time.Now(), *nFlag)

	// Allow the -from and -to flags to override the -n flag.
	var err error
	if len(*fromFlag) > 0 {
		cfg.DateFrom, err = time.ParseInLocation("2006-01-02", *fromFlag, time.Local)
		if err != nil {
			fail("bad format 'from' flag.")
		}
	}
	if len(*toFlag) > 0 {
		cfg.DateTo, err = time.ParseInLocation("2006-01-02", *toFlag, time.Local)
		if err != nil {
			fail("bad format 'to' flag.")
		}
	}

	tsSvc := svc.NewTimesheetSvc(cfg, svc.NewServices(cfg))
	tsSvc.Run()
}

// Month_offset returns the start and end date of the month
// that is n months before the origin date, o.
// So, if n == 0, it will return the start and end date of
// the month containing the origin date.
// Month_offset retains the Kind of its datetime argument.
func monthOffset(o time.Time, n int) (time.Time, time.Time) {
	dateFrom := time.Date(o.Year(), o.Month(), 1, 0, 0, 0, 0, o.Location()).AddDate(0, -n, 0)
	dateTo := time.Date(dateFrom.Year(), dateFrom.Month(), daysInMonth(dateFrom.Year(), dateFrom.Month()), 0, 0, 0, 0, o.Location())
	return dateFrom, dateTo
}

// DaysInMonth relies on the fact that time.Date allows month
// and day values to be normalized in their conversion to a date.
// We choose the -1th day of the next month to get the last day of
// the previous month.
func daysInMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

func fail(msg string) {
	fmt.Println(msg)
	flag.PrintDefaults()
	os.Exit(2)
}
