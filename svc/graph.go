package svc

import (
	"context"
	"fmt"
	"regexp"
	"time"

	azidentity "github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	msgraphsdk "github.com/microsoftgraph/msgraph-sdk-go"
	"github.com/microsoftgraph/msgraph-sdk-go/models"
	"github.com/microsoftgraph/msgraph-sdk-go/users"
	"github.com/vextasy/Timesheet_go/domain"
)

// graphSvc implements domain.GraphSvc.
type graphSvc struct {
	auth   domain.Auth
	client *msgraphsdk.GraphServiceClient
}

func NewGraphSvc(auth domain.Auth) domain.GraphSvc {
	cred, err := azidentity.NewClientSecretCredential(
		auth.TenantId,
		auth.ClientId,
		auth.ClientSecret,
		nil,
	)
	if err != nil {
		fmt.Printf("Error creating credentials: %v\n", err)
		return graphSvc{auth: auth, client: nil}
	}
	scopes := []string{"https://graph.microsoft.com/.default"}
	client, err := msgraphsdk.NewGraphServiceClientWithCredentials(cred, scopes)
	if err != nil {
		fmt.Printf("Error creating graph client: %v\n", err)
		return graphSvc{auth: auth, client: nil}
	}
	return graphSvc{
		auth:   auth,
		client: client,
	}
}

func (svc graphSvc) Read(userName string, fromDate time.Time, toDate time.Time) ([]domain.Task, error) {

	if svc.client == nil {
		return []domain.Task{}, nil
	}
	allusers, err := svc.client.Users().Get(context.Background(), nil)
	if err != nil {
		return []domain.Task{}, err
	}
	if allusers == nil || allusers.GetValue() == nil {
		return []domain.Task{}, nil
	}
	var targetUser models.Userable
	for _, user := range allusers.GetValue() {
		if *user.GetUserPrincipalName() == userName {
			targetUser = user
		}
	}
	if targetUser == nil {
		return []domain.Task{}, nil
	}

	// Got the user. Now get the tasks for that user.
	// Microsoft graph stores datetimes in UTC. So convert our range to UTC before filtering.
	// TODO: This doesn't actually change the dates it just changes the associated time zone.
	start := fromDate.UTC().Format("2006-01-02T15:04:05.0000000")
	end := toDate.UTC().Format("2006-01-02T15:04:05.0000000")
	filter := fmt.Sprintf("start/DateTime ge '%s' and start/DateTime le '%s' and IsAllDay eq false", start, end)
	query := users.ItemCalendarEventsRequestBuilderGetQueryParameters{
		Select: []string{"subject", "start", "end"},
		Filter: &filter,
		Top:    &[]int32{999}[0],
	}
	options := users.ItemCalendarEventsRequestBuilderGetRequestConfiguration{
		QueryParameters: &query,
	}
	events, err := svc.client.Users().ByUserId(*targetUser.GetId()).Calendar().Events().Get(context.Background(), &options)
	if err != nil {
		return []domain.Task{}, err
	}
	if events == nil || events.GetValue() == nil {
		return []domain.Task{}, nil
	}
	tasks := []domain.Task{}
	for _, ev := range events.GetValue() {
		if ev == nil || ev.GetSubject() == nil {
			continue
		}

		// proj (- group) - description
		var pat = `^\s*([\w/]+)(?:\s*-\s*([\w/]+))?\s*-\s*(.*)`

		matches := regexp.MustCompile(pat).FindStringSubmatch(*ev.GetSubject())
		if matches == nil || len(matches) != 4 { // Entire expression plus each subexpression.
			continue
		}
		proj := matches[1]
		group := matches[2]
		desc := matches[3]
		// Convert back from UTC to local time.
		// TODO: This doesn't actually change the dates it just changes the associated time zone.
		_start := ev.GetStart()
		_end := ev.GetEnd()
		start, err := time.Parse("2006-01-02T15:04:05.0000000", *_start.GetDateTime())
		if err != nil {
			return []domain.Task{}, fmt.Errorf("failed to parse start time: %v", err)
		}
		end, err := time.Parse("2006-01-02T15:04:05.0000000", *_end.GetDateTime())
		if err != nil {
			return []domain.Task{}, fmt.Errorf("failed to parse end time: %v", err)
		}
		start = start.Local()
		end = end.Local()
		duration := end.Sub(start)
		tasks = append(tasks, domain.Task{
			Project:  proj,
			Group:    group,
			Desc:     desc,
			Start:    start,
			Duration: duration,
		})
	}
	return tasks, nil
}
