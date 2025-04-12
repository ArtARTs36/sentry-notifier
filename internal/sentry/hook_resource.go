package sentry

import (
	"fmt"
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
