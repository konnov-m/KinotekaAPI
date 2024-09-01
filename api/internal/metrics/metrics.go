package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var NameMetric = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "names",
	Name:      "posts_methods",
	Help:      "Posts methods count",
})

func AddPosts() {
	NameMetric.Inc()
}
