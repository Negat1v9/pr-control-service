package middleware

import (
	"net/http"

	"github.com/Negat1v9/pr-review-service/pkg/metrics"
)

type Middleware func(http.Handler) http.Handler

type MiddleWareManager struct {
	metrics *metrics.PrometheusMetrics
}

func New(m *metrics.PrometheusMetrics) *MiddleWareManager {
	return &MiddleWareManager{
		metrics: m,
	}
}

func CreateStack(xs ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := 0; i < len(xs); i++ {
			x := xs[i]
			next = x(next)
		}
		return next
	}
}

// basic cors realisation
func CORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Headers", "Content-Type, Origin, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, Cache-Control, X-Requested-With")
		w.Header().Add("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")

		if r.Method == "OPTIONS" {
			http.Error(w, "No Content", http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
