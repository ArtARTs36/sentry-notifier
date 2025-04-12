package metrics

import (
	goMetrics "github.com/artarts36/go-metrics"
	"github.com/prometheus/client_golang/prometheus"
)

type Notifier struct {
	asyncQueueSize     *GaugeObserver
	asyncQueueCapacity prometheus.Gauge
}

func NewNotifier(registry goMetrics.Registry) *Notifier {
	return &Notifier{
		asyncQueueSize: newGaugeObserver(registry.NewGauge(prometheus.GaugeOpts{
			Name: "notifier_async_queue_size",
			Help: "Notifier (async): size of queue messages",
		})),
		asyncQueueCapacity: registry.NewGauge(prometheus.GaugeOpts{
			Name: "notifier_async_queue_capacity",
			Help: "Notifier (async): capacity of queue messages",
		}),
	}
}

func (n *Notifier) SetAsyncQueueCapacity(capacity int) {
	n.asyncQueueCapacity.Set(float64(capacity))
}

func (n *Notifier) ObserveAsyncQueueSize(callback func() int) {
	n.asyncQueueSize.callback = func() float64 {
		return float64(callback())
	}
}
