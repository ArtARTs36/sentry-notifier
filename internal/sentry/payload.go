package sentry

import (
	"fmt"
)

type Payload interface {
	GetHookResource() HookResource
	GetData() interface{}
}

func ParsePayload(resource HookResource, pl []byte) (Payload, error) {
	switch resource { //nolint: gocritic // cases will be append
	case HookResourceEventAlert:
		return createEventAlertFromJSON(pl)
	}

	return nil, fmt.Errorf("unknown resource: %s", resource)
}
