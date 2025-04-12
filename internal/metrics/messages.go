package metrics

import (
	goMetrics "github.com/artarts36/go-metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Messages struct {
	sentTotal    *prometheus.CounterVec
	sendingTotal *prometheus.CounterVec
}

func NewMessages(registry goMetrics.Registry) *Messages {
	return &Messages{
		sentTotal: registry.NewCounterVec(prometheus.CounterOpts{
			Name: "messages_sent_total",
			Help: "Count of sent messages",
		}, []string{"event", "messenger", "channel", "status"}),
		sendingTotal: registry.NewCounterVec(prometheus.CounterOpts{
			Name: "messages_sending_total",
			Help: "Count of sending messages",
		}, []string{"event", "messenger", "channel"}),
	}
}

func (m *Messages) IncSent(event, messenger, channel string, status bool) {
	st := "failed"
	if status {
		st = "success"
	}

	m.sentTotal.WithLabelValues(event, messenger, channel, st).Inc()
}

func (m *Messages) IncSending(event, messenger, channel string) {
	m.sendingTotal.WithLabelValues(event, messenger, channel).Inc()
}
