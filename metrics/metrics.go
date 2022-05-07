package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	label             = []string{"app"}
	NackTotalMessages = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "nack_messages_total",
		Namespace: "nats",
	}, label)

	PublishErrorCount = promauto.NewCounterVec(prometheus.CounterOpts{
		Name:      "publish_error_total",
		Namespace: "nats",
	}, label)
)
