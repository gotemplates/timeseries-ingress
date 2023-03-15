package handler

import (
	"github.com/gotemplates/core/exchange"
	"github.com/gotemplates/core/runtime"
	"github.com/gotemplates/host/messaging"
	"net/http"
)

func SetHealthPatterns(host, postgres string) {
	hostPattern = host
	postgresPattern = postgres
}

func SetPgxsqlUri(uri string) {
	pgxsqlUri = uri
}

var (
	hostPattern     string
	postgresPattern string
	pgxsqlUri       string
)

func HealthLivenessHandler(w http.ResponseWriter, r *http.Request) {
	var status = runtime.NewHttpStatusCode(http.StatusBadRequest)
	if r == nil || r.URL == nil {
		exchange.WriteResponse(w, nil, status)
		return
	}
	if r.URL.Path == hostPattern {
		status = runtime.NewStatusOK()
	} else {
		if r.URL.Path == postgresPattern {
			status = messaging.Ping[runtime.LogError](r.Context(), pgxsqlUri)
		}
	}
	if status.OK() {
		exchange.WriteResponse(w, []byte("up"), status)
	} else {
		exchange.WriteResponse(w, nil, status)
	}
}
