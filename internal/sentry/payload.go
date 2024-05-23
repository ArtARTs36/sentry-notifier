package sentry

import (
	"fmt"
)

type Payload interface {
	GetHookResource() HookResource
	GetData() interface{}
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
