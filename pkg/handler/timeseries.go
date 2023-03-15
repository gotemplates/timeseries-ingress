package handler

import (
	"github.com/go-http-utils/headers"
	"github.com/gotemplates/core/exchange"
	"github.com/gotemplates/core/runtime"
	"github.com/gotemplates/timeseries/accesslog"
	"net/http"
)

func GetContentLocation(req *http.Request) string {
	if req != nil && req.Header != nil {
		return req.Header.Get(headers.ContentLocation)
	}
	return ""
}

func TimeseriesHandler(w http.ResponseWriter, r *http.Request) {
	if r == nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	switch r.Method {
	case "GET":
		data, status := accesslog.GetByte[runtime.LogError](runtime.ContextWithRequest(r), GetContentLocation(r), r.URL.Query())
		exchange.WriteResponse(w, data, status, headers.ContentType)
	case "PUT":
		data, err := exchange.ReadAll(r.Body)
		if err != nil {
			exchange.WriteResponse(w, nil, runtime.Handle[runtime.LogError]()("/timeseries/handler", err))
			return
		}
		_, status := accesslog.PutByte[runtime.LogError](runtime.ContextWithRequest(r), GetContentLocation(r), data)
		exchange.WriteResponse(w, nil, status)
	default:
		w.WriteHeader(http.StatusBadRequest)
	}
}
