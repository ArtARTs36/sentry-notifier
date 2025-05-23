package sentry

import (
	"fmt"
)

type Payload interface {
	GetID() string
	GetHookResource() HookResource
	GetData() interface{}
	GetProjectSlug() string
}

func ParsePayload(resource HookResource, pl []byte) (Payload, error) {
	switch resource {
	case HookResourceEventAlert:
		return createEventAlertFromJSON(pl)
	case HookResourceIssue:
		return createIssuePayload(pl)
	}

	return nil, fmt.Errorf("unknown resource: %s", resource)
}
