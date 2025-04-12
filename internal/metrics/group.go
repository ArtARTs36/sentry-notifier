package metrics

import goMetrics "github.com/artarts36/go-metrics"

type Group struct {
	Messages *Messages
}

func NewGroup(registry goMetrics.Registry) *Group {
	return &Group{
		Messages: NewMessages(registry),
	}
}
