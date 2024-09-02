package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var MetricPosts = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "names",
	Name:      "posts_methods",
	Help:      "Posts methods count",
})

var MetricGet = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "names",
	Name:      "get_methods",
	Help:      "Get methods count",
})

var AllMetric = promauto.NewCounter(prometheus.CounterOpts{
	Namespace: "names",
	Name:      "all_methods",
	Help:      "All methods count",
})

func AddPosts() {
	MetricPosts.Inc()
}

func AddGet() {
	MetricGet.Inc()
}

func AddAll() {
	AllMetric.Inc()
}
