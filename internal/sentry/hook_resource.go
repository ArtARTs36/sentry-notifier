package sentry

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type HookResource string

const (
	HookResourceEventAlert = "event_alert"
	HookResourceIssue      = "issue"
)

var HookResourceMap = map[string]HookResource{
	HookResourceEventAlert: HookResourceEventAlert,
	HookResourceIssue:      HookResourceIssue,
}

func WrapHookResource(resource string) (HookResource, error) {
	res, ok := HookResourceMap[resource]
	if !ok {
		return "", fmt.Errorf("unknown hook resource %q", resource)
	}

	return res, nil
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
