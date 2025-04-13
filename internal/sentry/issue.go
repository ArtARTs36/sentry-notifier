package sentry

import (
	"encoding/json"
	"github.com/google/uuid"
	"time"
)

type IssuePayload struct {
	Action string `json:"action"`
	Data   struct {
		Issue struct {
			Action    string `json:"-"`
			Count     string `json:"count"`
			ID        string `json:"id"`
			Level     string `json:"level"`
			ShortID   string `json:"shortId"`
			Status    string `json:"status"`
			Type      string `json:"type"`
			Title     string `json:"title"`
			LastSeen  Time   `json:"lastSeen"`
			FirstSeen Time   `json:"firstSeen"`
			Project   struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Platform string `json:"platform"`
				Slug     string `json:"slug"`
			} `json:"project"`
		} `json:"issue"`
	} `json:"data"`
}

func (p *IssuePayload) GetProjectSlug() string {
	return p.Data.Issue.Project.Slug
}

func createIssuePayload(data []byte) (*IssuePayload, error) {
	pl := new(IssuePayload)
	err := json.Unmarshal(data, &pl)
	if err != nil {
		return nil, err
	}

	pl.Data.Issue.Action = pl.Action

	return pl, nil
}

func (p *IssuePayload) GetID() string {
	return p.Data.Issue.ID
}

func (p *IssuePayload) GetHookResource() HookResource {
	return HookResourceIssue
}

func (p *IssuePayload) GetData() interface{} {
	return p.Data
}

func ExampleIssuePayload() *IssuePayload {
	pl := &IssuePayload{
		Action: "test",
	}
	pl.Data.Issue.Action = "test"
	pl.Data.Issue.ID = uuid.NewString()
	pl.Data.Issue.Level = "error"
	pl.Data.Issue.ShortID = pl.Data.Issue.ID
	pl.Data.Issue.Status = "unknown"
	pl.Data.Issue.Type = "unknown"
	pl.Data.Issue.Title = "Test Issue"
	pl.Data.Issue.FirstSeen.Time = time.Now()
	pl.Data.Issue.LastSeen.Time = time.Now()
	pl.Data.Issue.Project.ID = uuid.NewString()
	pl.Data.Issue.Project.Name = "test-project"
	pl.Data.Issue.Project.Platform = "test-platform"
	pl.Data.Issue.Project.Slug = "test_project"

	return pl
}
