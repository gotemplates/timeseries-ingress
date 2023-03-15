package host

import (
	"errors"
	"github.com/gotemplates/core/exchange"
	"github.com/gotemplates/core/runtime"
	"github.com/gotemplates/host/accessdata"
	"github.com/gotemplates/host/accesslog"
	"github.com/gotemplates/host/controller"
	"github.com/gotemplates/host/messaging"
	middleware2 "github.com/gotemplates/host/middleware"
	runtime2 "github.com/gotemplates/host/runtime"
	"github.com/gotemplates/timeseries-ingress/pkg/resource"
	"net/http"
	"time"
)

const (
	startupLoc                = "/host/startup"
	egressLogOperatorNameFmt  = "fs/egress_logging_operators_{env}.json"
	ingressLogOperatorNameFmt = "fs/ingress_logging_operators_{env}.json"
)

func Startup[E runtime.ErrorHandler, O runtime.OutputHandler](mux *http.ServeMux) (http.Handler, *runtime.Status) {
	var e E

	initOrigin()
	err := initLogging()
	if err != nil {
		return nil, e.Handle(startupLoc+"/logging", err)
	}
	errs := initControllers()
	if len(errs) > 0 {
		return nil, e.Handle(startupLoc+"/controllers", errs...)
	}
	initMux(mux)
	status := startupResources[E, O](true)
	if !status.OK() {
		return mux, status
	}

	middleware2.ControllerWrapTransport(exchange.Client)
	return middleware2.ControllerHttpHostMetricsHandler(mux, ""), status
}

func Shutdown() {
	messaging.Shutdown()
}

func startupResources[E runtime.ErrorHandler, O runtime.OutputHandler](addCredentials bool) *runtime.Status {
	var e E
	content, err := createContent(addCredentials)
	if err != nil {
		return e.Handle(startupLoc+"/resources", err)
	}
	return messaging.Startup[E, O](time.Second*5, content)

}

func createContent(addCredentials bool) (messaging.ContentMap, error) {
	content := make(messaging.ContentMap, 2)
	if config.DatabaseUrl() == "" {
		return nil, errors.New("database url is empty")
	}
	content[config.PostgresPgxsqlUri()] = []any{messaging.DatabaseUrl{Url: config.DatabaseUrl()}, messaging.ControllerApply(controller.EgressApply)}
	if addCredentials {
		content[config.PostgresPgxsqlUri()] = append(content[config.PostgresPgxsqlUri()], messaging.Credentials(func() (username string, password string, err error) { return "", "", nil }))
	}
	return content, nil
}

func initOrigin() {
	/*	shared.SetOrigin(shared.Origin{
			Region:     "Region",
			Zone:       "Zone",
			SubZone:    "SubZone",
			Service:    "Service",
			InstanceId: "InstanceId",
		})

	*/
}

func initLogging() error {
	// Options that are defaulted to true for the statuses
	accesslog.SetIngressLogStatus(true)
	accesslog.SetEgressLogStatus(true)
	accesslog.SetPingLogStatus(true)

	// Enable logging function for access events middleware
	// middleware.SetLogFn(func(entry *data.Entry) {
	// log.Write[log.DebugOutputHandler, data.JsonFormatter](entry)
	//},
	//)
	controller.SetLogFn(func(traffic string, start time.Time, duration time.Duration, req *http.Request, resp *http.Response, statusFlags string, controllerState map[string]string) {
		entry := accessdata.NewEntry(traffic, start, duration, req, resp, statusFlags, controllerState)
		accesslog.Write[accesslog.DebugOutputHandler, accessdata.JsonFormatter](entry)
	},
	)

	err := accesslog.CreateIngressOperators(func() ([]byte, error) {
		return resource.ReadFile(runtime2.EnvExpansion(ingressLogOperatorNameFmt))
	})
	if err == nil {
		err = accesslog.CreateEgressOperators(func() ([]byte, error) {
			return resource.ReadFile(runtime2.EnvExpansion(egressLogOperatorNameFmt))
		})
	}
	return err
}
