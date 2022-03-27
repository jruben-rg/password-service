package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/jruben-rg/password-service/go-pwned/handlers"
	"github.com/jruben-rg/password-service/go-pwned/metric"
	"github.com/jruben-rg/password-service/go-pwned/middleware"
	"github.com/jruben-rg/password-service/go-pwned/password"
	"github.com/jruben-rg/password-service/go-pwned/pwned"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	// Change logging to stdout
	log.SetOutput(os.Stdout)
}

func main() {

	log := log.New(os.Stdout, "goPwned", log.LstdFlags)

	if len(os.Args) < 2 {
		log.Fatal("provide a path for a yaml file so configuration can be loaded")
	}

	//Initialise pwned and password validators
	passwordValidator := password.NewPasswordConfig(os.Args[1])
	pwnedValidator := pwned.NewPwnedConfig(os.Args[1])

	metricsService, err := metric.NewPrometheusService()
	if err != nil {
		log.Fatal(err)
	}

	//Chain handlers
	pwnedHandler := handlers.NewPwnedHandler(log, pwnedValidator.IsSecurePassword)
	passwordHandler := handlers.NewPasswordHandler(log, passwordValidator.Validate, pwnedHandler)
	healthtzHandler := handlers.NewHealthzHandler(log)

	mux := http.NewServeMux()
	mux.Handle("/validate", passwordHandler)
	mux.Handle("/healthz", healthtzHandler)
	mux.Handle("/metrics", promhttp.Handler())
	wrappedMux := middleware.Metrics(metricsService, mux)

	s := &http.Server{
		Addr:              ":2112",
		Handler:           wrappedMux,
		ReadTimeout:       1 * time.Second,
		ReadHeaderTimeout: 0,
		WriteTimeout:      1 * time.Second,
		IdleTimeout:       120 * time.Second,
		MaxHeaderBytes:    0,
	}

	go func() {
		err := s.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <-sigChan
	log.Println("Received Terminate, gracefully shutdown", sig)

	tc, cancelFunc := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelFunc()
	s.Shutdown(tc)
}
