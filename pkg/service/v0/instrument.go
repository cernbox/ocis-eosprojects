package service

import (
	"github.com/cernbox/ocis-eosprojects/pkg/metrics"
	"github.com/cernbox/ocis-eosprojects/pkg/proto/v0"
	userpb "github.com/cs3org/go-cs3apis/cs3/identity/user/v1beta1"
	"github.com/prometheus/client_golang/prometheus"
)

// NewInstrument returns a service that instruments metrics.
func NewInstrument(next EosProjects, metrics *metrics.Metrics) EosProjects {
	return instrument{
		next:    next,
		metrics: metrics,
	}
}

type instrument struct {
	next    EosProjects
	metrics *metrics.Metrics
}

// Greet implements the EosProjects interface.
func (i instrument) GetProjects(user *userpb.User) []*proto.Project {
	timer := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		us := v * 1000000

		i.metrics.Latency.WithLabelValues().Observe(us)
		i.metrics.Duration.WithLabelValues().Observe(v)
	}))

	defer timer.ObserveDuration()

	i.metrics.Counter.WithLabelValues().Inc()

	return i.next.GetProjects(user)
}
