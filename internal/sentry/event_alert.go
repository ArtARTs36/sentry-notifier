package sentry

import (
	"encoding/json"
	"log/slog"
	"net/url"
	"slices"
	"strings"
)

type EventAlert struct {
	Action string `json:"action"`
	Data   struct {
		Event struct {
			IssueURL string `json:"issue_url"`
			IssueID  string `json:"issue_id"`
			Platform string `json:"platform"`
			Title    string `json:"title"`
			Type     string `json:"type"`
			Project  int    `json:"project"`
			URL      string `json:"url"`
			Datetime Time   `json:"datetime"`
			WebURL   string `json:"web_url"`
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
