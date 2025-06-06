package metrics

import (
	"github.com/prometheus/client_golang/prometheus"

	goMetrics "github.com/artarts36/go-metrics"
)

type AppInfo struct {
	info *prometheus.GaugeVec
}

func NewAppInfo(registry goMetrics.Registry) *AppInfo {
	return &AppInfo{
		info: registry.NewGaugeVec(prometheus.GaugeOpts{
			Name: "app_info",
			Help: "App info",
		}, []string{"version", "build_time", "supported_messengers"}),
	}
}

func (i *AppInfo) SetInfo(version string, buildTime string, supportedMessengers string) {
	i.info.WithLabelValues(version, buildTime, supportedMessengers)
}
