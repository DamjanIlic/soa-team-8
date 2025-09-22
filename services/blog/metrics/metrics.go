package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

var (
	BlogPostsProcessed = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "blog_posts_processed_total",
			Help: "Ukupan broj obrađenih blog postova",
		},
	)
)

func Init() {
	prometheus.MustRegister(BlogPostsProcessed)
}
