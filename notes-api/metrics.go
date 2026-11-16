package main

// Kennzahlen mit der offiziellen Bibliothek client_golang:
// https://prometheus.io/docs/guides/go-application/

import (
	"net/http"
	"strconv"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/trace"
)

type metrics struct {
	reg      *prometheus.Registry
	requests *prometheus.CounterVec
	duration *prometheus.HistogramVec
	created  prometheus.Counter
	dbUp     prometheus.Gauge
}

func newMetrics(store string) *metrics {
	reg := prometheus.NewRegistry()
	m := &metrics{
		reg: reg,
		requests: prometheus.NewCounterVec(prometheus.CounterOpts{
			Name: "notes_http_requests_total",
			Help: "HTTP-Anfragen nach Methode, Pfad und Status.",
		}, []string{"method", "path", "code"}),
		duration: prometheus.NewHistogramVec(prometheus.HistogramOpts{
			Name:    "notes_http_request_duration_seconds",
			Help:    "Dauer der HTTP-Anfragen in Sekunden.",
			Buckets: []float64{.005, .01, .025, .05, .1, .25, .5, 1, 2.5},
		}, []string{"method", "path"}),
		created: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "notes_created_total",
			Help: "Angelegte Notizen seit dem Start.",
		}),
		dbUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "notes_db_up",
			Help: "1, wenn der letzte Healthcheck die Datenbank erreicht hat.",
		}),
	}
	info := prometheus.NewGauge(prometheus.GaugeOpts{
		Name:        "notes_build_info",
		Help:        "Version und Speicherart von notes-api.",
		ConstLabels: prometheus.Labels{"version": version, "store": store},
	})
	info.Set(1)
	reg.MustRegister(m.requests, m.duration, m.created, info,
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}))
	if store == "postgres" {
		reg.MustRegister(m.dbUp)
	}
	return m
}

func (m *metrics) noteCreated() { m.created.Inc() }

func (m *metrics) setDBUp(up bool) {
	if up {
		m.dbUp.Set(1)
	} else {
		m.dbUp.Set(0)
	}
}

// handler liefert /metrics – mit OpenMetrics, sonst fehlen die Exemplars.
func (m *metrics) handler() http.Handler {
	return promhttp.HandlerFor(m.reg, promhttp.HandlerOpts{EnableOpenMetrics: true})
}

// Nur bekannte Pfade als Label – sonst wächst die Zahl der Zeitreihen unbegrenzt.
var knownPaths = map[string]bool{
	"/healthz": true, "/version": true, "/notes": true, "/upstream": true, "/metrics": true,
}

type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (r *statusRecorder) WriteHeader(code int) {
	r.code = code
	r.ResponseWriter.WriteHeader(code)
}

// middleware zählt, misst und schreibt eine Logzeile pro Anfrage.
func (m *metrics) middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, code: http.StatusOK}
		next.ServeHTTP(rec, r)
		took := time.Since(start)

		path := r.URL.Path
		if !knownPaths[path] {
			path = "other"
		}
		m.requests.WithLabelValues(r.Method, path, strconv.Itoa(rec.code)).Inc()
		obs := m.duration.WithLabelValues(r.Method, path)
		// Exemplar: die Trace-ID dieser Anfrage hängt am Messwert (Tag 11).
		if sc := trace.SpanContextFromContext(r.Context()); sc.IsSampled() {
			obs.(prometheus.ExemplarObserver).ObserveWithExemplar(took.Seconds(),
				prometheus.Labels{"trace_id": sc.TraceID().String()})
		} else {
			obs.Observe(took.Seconds())
		}
		if path == "/metrics" || path == "/healthz" {
			return // Abfragen von Prometheus und Kubernetes nicht loggen
		}
		logRequest(r, rec.code, took)
	})
}
