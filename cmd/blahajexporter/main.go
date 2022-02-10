package main

import (
	"flag"
	"fmt"
	"github.com/patrick246/blahaj-exporter/internal/client"
	"github.com/patrick246/blahaj-exporter/internal/exporter"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"

	"go.uber.org/zap"

	// This controls the maxprocs environment variable in container runtimes.
	// see https://martin.baillie.id/wrote/gotchas-in-the-go-network-packages-defaults/#bonus-gomaxprocs-containers-and-the-cfs
	_ "go.uber.org/automaxprocs"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "an error occurred: %s\n", err)
		os.Exit(1)
	}
}

func run() error {
	addr := flag.String("addr", ":8080", "Port to listen on")
	contact := flag.String("contact", "", "Contact to put in the user agent")
	flag.Parse()

	logger, err := zap.NewProduction()
	if err != nil {
		return err
	}

	defer logger.Sync()

	log := logger.Sugar()

	if *contact == "" {
		log.Fatalw("mandatory flag missing", "flag", "contact")
	}

	stockClient := client.New(&http.Client{}, *contact)
	collector := exporter.New(stockClient, logger.Sugar())

	err = prometheus.Register(collector)
	if err != nil {
		return err
	}

	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())

	srv := http.Server{
		Addr:    *addr,
		Handler: mux,
	}

	log.Infow("listening", "addr", addr)
	err = srv.ListenAndServe()
	if err != nil {
		return err
	}

	return nil
}
