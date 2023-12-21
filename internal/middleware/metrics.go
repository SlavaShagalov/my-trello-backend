package middleware

import (
	"github.com/SlavaShagalov/my-trello-backend/internal/pkg/metrics"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/urfave/negroni"
)

func NewMetrics(mt metrics.PrometheusMetrics) func(h http.HandlerFunc) http.HandlerFunc {
	return func(h http.HandlerFunc) http.HandlerFunc {
		return func(w http.ResponseWriter, r *http.Request) {
			wWithCode := negroni.NewResponseWriter(w)

			begin := time.Now()
			h(wWithCode, r)
			httpCode := wWithCode.Status()
			mt.ExecutionTime().
				WithLabelValues(strconv.Itoa(httpCode), r.URL.String(), r.Method).
				Observe(float64(time.Since(begin).Milliseconds()))

			log.Println(strconv.Itoa(httpCode), r.URL.String(), r.Method, "MS", time.Since(begin).Milliseconds())

			mt.TotalHits().Inc()

			if 200 <= httpCode && httpCode <= 399 {
				mt.SuccessHits().
					WithLabelValues(strconv.Itoa(httpCode), r.URL.String(), r.Method).Inc()
			} else {
				mt.ErrorsHits().
					WithLabelValues(strconv.Itoa(httpCode), r.URL.String(), r.Method).Inc()
			}
		}
	}
}
