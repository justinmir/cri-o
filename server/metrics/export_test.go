package metrics

import "github.com/prometheus/client_golang/prometheus"

func (m *Metrics) ImagePullsSuccessByRegistryVec() *prometheus.CounterVec {
	return m.metricImagePullsSuccessByRegistryTotal
}

func (m *Metrics) ImagePullsFailureByRegistryVec() *prometheus.CounterVec {
	return m.metricImagePullsFailureByRegistryTotal
}
