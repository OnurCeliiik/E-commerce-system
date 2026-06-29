package metrics

import "github.com/prometheus/client_golang/prometheus"

var InventoryEventsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "inventory_events_total",
		Help: "order.created events handled by inventory-service",
	},
	[]string{"event", "result"},
)

func init() {
	prometheus.MustRegister(InventoryEventsTotal)
}

func RecordInventoryEvent(event, result string) {
	InventoryEventsTotal.WithLabelValues(event, result).Inc()
}
