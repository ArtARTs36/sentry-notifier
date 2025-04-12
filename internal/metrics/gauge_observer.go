package metrics

import "github.com/prometheus/client_golang/prometheus"

type GaugeObserver struct {
	gauge    prometheus.Gauge
	callback func() float64
}

func newGaugeObserver(gauge prometheus.Gauge) *GaugeObserver {
	return &GaugeObserver{
		gauge: gauge,
		callback: func() float64 {
			return 0
		},
	}
}

func (o *GaugeObserver) Describe(desc chan<- *prometheus.Desc) {
	o.gauge.Describe(desc)
}

func (o *GaugeObserver) Collect(metric chan<- prometheus.Metric) {
	o.gauge.Set(o.callback())
	o.gauge.Collect(metric)
}
