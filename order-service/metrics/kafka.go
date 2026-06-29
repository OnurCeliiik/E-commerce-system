package metrics

import "github.com/prometheus/client_golang/prometheus"

var InventoryEventsTotal = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "order_inventory_events_total",
		Help: "Inventory outcome events processed by order-service",
	},
	[]string{"event", "result"},
)

func init() {
	prometheus.MustRegister(InventoryEventsTotal)
}

func RecordInventoryEvent(event, result string) {
	InventoryEventsTotal.WithLabelValues(event, result).Inc()
}
