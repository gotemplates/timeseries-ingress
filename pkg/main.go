package main

import (
	"context"
	"fmt"
	runtime2 "github.com/gotemplates/core/runtime"
	"github.com/gotemplates/timeseries-ingress/pkg/host"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"
)

const (
	addr         = "0.0.0.0:8080"
	writeTimeout = time.Second * 300
	readTimeout  = time.Second * 15
	idleTimeout  = time.Second * 60
)

func main() {
	start := time.Now()
	displayRuntime[runtime2.StdOutput]()
	handler, status := host.Startup[runtime2.LogError, runtime2.StdOutput](http.NewServeMux())
	if !status.OK() {
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("host started: %v", time.Since(start)))
	defer host.Shutdown()

	srv := http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: writeTimeout,
		ReadTimeout:  readTimeout,
		IdleTimeout:  idleTimeout,
		Handler:      handler,
	}

	idleConnsClosed := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt)
		<-sigint

		// We received an interrupt signal, shut down.
		if err := srv.Shutdown(context.Background()); err != nil {
			// Error from closing listeners, or context timeout:
			log.Printf("HTTP server Shutdown: %v", err)
		} else {
			log.Printf("HTTP server Shutdown")
		}
		close(idleConnsClosed)
	}()

	if err := srv.ListenAndServe(); err != http.ErrServerClosed {
		// Error starting or closing listener:
		log.Fatalf("HTTP server ListenAndServe: %v", err)
	}

	<-idleConnsClosed
}

func displayRuntime[O runtime2.OutputHandler]() {
	var o O
	o.Write(fmt.Sprintf("vers : %v", runtime.Version()))
	o.Write(fmt.Sprintf("os   : %v", runtime.GOOS))
	o.Write(fmt.Sprintf("arch : %v", runtime.GOARCH))
	o.Write(fmt.Sprintf("cpu  : %v", runtime.NumCPU()))
}
