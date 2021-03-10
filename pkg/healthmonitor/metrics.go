package healthmonitor

import (
	compbasemetrics "k8s.io/component-base/metrics"
)

type registerables []compbasemetrics.Registerable

var (
	healthyTargetsTotal = compbasemetrics.NewCounterVec(
		&compbasemetrics.CounterOpts{
			Name:           "health_monitor_healthy_target_total",
			Help:           "Number of healthy instances registered with the health monitor. Partitioned by targets",
			StabilityLevel: compbasemetrics.ALPHA,
		},
		[]string{"target"},
	)

	currentHealthyTargets = compbasemetrics.NewGauge(
		&compbasemetrics.GaugeOpts{
			Name:           "health_monitor_current_healthy_targets",
			Help:           "Number of currently healthy instances observed by the health monitor",
			StabilityLevel: compbasemetrics.ALPHA,
		},
	)

	unHealthyTargetsTotal = compbasemetrics.NewCounterVec(
		&compbasemetrics.CounterOpts{
			Name:           "health_monitor_unhealthy_target_total",
			Help:           "Number of unhealthy instances registered with the health monitor. Partitioned by targets",
			StabilityLevel: compbasemetrics.ALPHA,
		},
		[]string{"target"},
	)

	metrics = registerables{
		healthyTargetsTotal,
		currentHealthyTargets,
		unHealthyTargetsTotal,
	}
)

// HealthyTargetsTotal increments the total number of healthy instances observed by the health monitor
func HealthyTargetsTotal(target string) {
	healthyTargetsTotal.WithLabelValues(target).Add(1)
}

// CurrentHealthyTargets keeps track of the current number of healthy targets observed by the health monitor
func CurrentHealthyTargets(count float64) {
	currentHealthyTargets.Set(count)
}

// UnHealthyTargetsTotal increments the total number of unhealthy instances observed by the health monitor
func UnHealthyTargetsTotal(target string) {
	unHealthyTargetsTotal.WithLabelValues(target).Add(1)
}

// Metrics specifies a set of methods that are used to register various metrics
type Metrics struct {
	// HealthyTargetsTotal increments the total number of healthy instances observed by the health monitor
	HealthyTargetsTotal func(target string)

	// CurrentHealthyTargets keeps track of the current number of healthy targets observed by the health monitor
	CurrentHealthyTargets func(count float64)

	// UnHealthyTargetsTotal increments the total number of unhealthy instances observed by the health monitor
	UnHealthyTargetsTotal func(target string)
}

// Register is a way to register the health monitor related metrics in the provided store
func Register(registerFn func(...compbasemetrics.Registerable)) *Metrics {
	registerFn(metrics...)
	return &Metrics{
		HealthyTargetsTotal:   HealthyTargetsTotal,
		CurrentHealthyTargets: CurrentHealthyTargets,
		UnHealthyTargetsTotal: UnHealthyTargetsTotal,
	}
}

type noopMetrics struct{}

func (noopMetrics) TargetsTotal(string)  {}
func (noopMetrics) TargetsGauge(float64) {}
