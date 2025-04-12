package sentry

import (
	"encoding/json"
	"github.com/google/uuid"
	"log/slog"
	"net/url"
	"slices"
	"strings"
	"time"
)

type EventAlert struct {
	Action string `json:"action"`
	Data   struct {
		Event struct {
			IssueURL    string   `json:"issue_url"`
			IssueID     string   `json:"issue_id"`
			Platform    string   `json:"platform"`
			Title       string   `json:"title"`
			Type        string   `json:"type"`
			Project     int      `json:"project"`
			URL         string   `json:"url"`
			Datetime    Time     `json:"datetime"`
			WebURL      string   `json:"web_url"`
			Fingerprint []string `json:"fingerprint"`
			Request     struct {
				Method string `json:"method"`
				URL    string `json:"url"`
			} `json:"request"`
		} `json:"event"`
		Extracted struct {
			ProjectName      string `json:"-"`
			OrganizationName string `json:"-"`
		} `json:"-"`
	} `json:"data"`
}

func createEventAlertFromJSON(data []byte) (*EventAlert, error) {
	a := new(EventAlert)
	err := json.Unmarshal(data, &a)
	if err != nil {
		return nil, err
	}

	a.extract()

	return a, nil
}

func (a *EventAlert) GetID() string {
	return a.Data.Event.IssueID
}

func (a *EventAlert) GetHookResource() HookResource {
	return HookResourceEventAlert
}

func (a *EventAlert) GetData() interface{} {
	return a.Data
}

func (a *EventAlert) extract() {
	ur, err := url.Parse(a.Data.Event.URL)
	if err != nil {
		slog.
			With(slog.String("err", err.Error())).
			Warn("failed to parse url")
	} else {
		urParts := strings.Split(ur.Path, "/")
		startIndex := slices.Index(urParts, "projects")
		if startIndex >= 0 && startIndex+1 <= len(urParts) {
			a.Data.Extracted.OrganizationName = urParts[startIndex+1]

			if startIndex+2 < len(urParts) {
				a.Data.Extracted.ProjectName = urParts[startIndex+2]
			}
		}
	}
}

func ExampleEventAlert() *EventAlert {
	pl := &EventAlert{
		Action: "test",
	}

	pl.Data.Event.IssueURL = "https://sentry.io/api/0/projects/test-org/front-end/events/e4874d664c3540c1a32eab185f12/"
	pl.Data.Event.IssueID = uuid.NewString()
	pl.Data.Event.Platform = "web"
	pl.Data.Event.Title = "Test Alert"
	pl.Data.Event.Type = "Test Type"
	pl.Data.Event.Project = 1
	pl.Data.Event.URL = "https://sentry.io/api/0/projects/test-org/front-end/events/e4874d664c3540c1a32eab185f12c5ab/"
	pl.Data.Event.Datetime.Time = time.Now()
	pl.Data.Event.WebURL = "https://sentry.io/api/0/projects/test-org/front-end/events/e4874d664c3540c1a32eab185f12c/"
	pl.Data.Event.Request.Method = "GET"
	pl.Data.Event.Request.URL = "/v1/users/"
	pl.Data.Extracted.ProjectName = "TestProject"
	pl.Data.Extracted.OrganizationName = "TestOrganization"

	return pl
}
