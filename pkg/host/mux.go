package host

import (
	"github.com/gotemplates/timeseries-ingress/pkg/handler"
	"net/http"
	"net/http/pprof"
)

const (
	accessLogPattern              = "/access-log"
	healthLivenessPattern         = "/health/liveness"
	healthLivenessPostgresPattern = "/health/liveness/postgres"

	IndexPattern   = "/debug/pprof/"
	CmdLinePattern = "/debug/pprof/cmdline"
	ProfilePattern = "/debug/pprof/profile" // ?seconds=30
	SymbolPattern  = "/debug/pprof/symbol"
	TracePattern   = "/debug/pprof/trace"
)

func initMux(r *http.ServeMux) {
	handler.SetHealthPatterns(healthLivenessPattern, healthLivenessPostgresPattern)
	handler.SetPgxsqlUri(config.PostgresPgxsqlUri())

	addProfileRoutes(r)
	r.Handle(accessLogPattern, http.HandlerFunc(handler.TimeseriesHandler))
	r.Handle(healthLivenessPattern, http.HandlerFunc(handler.HealthLivenessHandler))
	r.Handle(healthLivenessPostgresPattern, http.HandlerFunc(handler.HealthLivenessHandler))

}

func addProfileRoutes(r *http.ServeMux) {
	r.Handle(IndexPattern, http.HandlerFunc(pprof.Index))
	r.Handle(CmdLinePattern, http.HandlerFunc(pprof.Cmdline))
	r.Handle(ProfilePattern, http.HandlerFunc(pprof.Profile))
	r.Handle(SymbolPattern, http.HandlerFunc(pprof.Symbol))
	r.Handle(TracePattern, http.HandlerFunc(pprof.Trace))

}
