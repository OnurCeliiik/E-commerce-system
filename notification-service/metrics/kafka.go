package metrics

import "github.com/prometheus/client_golang/prometheus"

var (
	NotificationEventsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "notification_events_total",
			Help: "Notification events processed by notification-service",
		},
		[]string{"event", "result"},
	)
)

func init() {
	prometheus.MustRegister(NotificationEventsTotal)
}

func RecordNotificationEvent(event, result string) {
	NotificationEventsTotal.WithLabelValues(event, result).Inc()
}
