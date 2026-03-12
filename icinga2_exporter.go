package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/martialblog/icinga2-exporter/internal/collector"
	"github.com/martialblog/icinga2-exporter/internal/icinga"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// nolint: gochecknoglobals
var (
	// These get filled at build time with the proper vaules.
	version = "development"
	commit  = "HEAD"
	date    = "latest"
)

func buildVersion() string {
	result := version

	if commit != "" {
		result = fmt.Sprintf("%s\ncommit: %s", result, commit)
	}

	if date != "" {
		result = fmt.Sprintf("%s\ndate: %s", result, date)
	}

	return result
}

func main() {
	var (
		cliListenAddress        string
		cliMetricsPath          string
		cliCAFile               string
		cliCertFile             string
		cliKeyFile              string
		cliUsername             string
		cliPassword             string
		cliBaseURL              string
		cliCacheTTL             uint
		cliVersion              bool
		cliDebugLog             bool
		cliInsecure             bool
		cliCollectorApiListener bool
		cliCollectorCIB         bool
		cliCollectorChecker     bool
		cliCollectorInflux      bool
		cliCollectorInflux2     bool
		cliCollectorGraphite    bool
		cliCollectorOTLP        bool
	)

	flag.StringVar(&cliListenAddress, "web.listen-address", ":9665", "Address on which to expose metrics and web interface.")
	flag.StringVar(&cliMetricsPath, "web.metrics-path", "/metrics", "Path under which to expose metrics.")
	flag.UintVar(&cliCacheTTL, "web.cache-ttl", 60, "Cache lifetime in seconds for the Icinga API responses")

	flag.StringVar(&cliBaseURL, "icinga.api", "https://localhost:5665/v1", "Path to the Icinga2 API")
	flag.StringVar(&cliUsername, "icinga.username", "", "Icinga2 API Username")
	flag.StringVar(&cliPassword, "icinga.password", "", "Icinga2 API Password")
	flag.StringVar(&cliCAFile, "icinga.cafile", "", "Path to the Icinga2 API TLS CA")
	flag.StringVar(&cliCertFile, "icinga.certfile", "", "Path to the Icinga2 API TLS cert")
	flag.StringVar(&cliKeyFile, "icinga.keyfile", "", "Path to the Icinga2 API TLS key")
	flag.BoolVar(&cliInsecure, "icinga.insecure", false, "Skip TLS verification for Icinga2 API")

	flag.BoolVar(&cliCollectorApiListener, "collector.apilistener", false, "Include APIListener data")
	flag.BoolVar(&cliCollectorCIB, "collector.cib", false, "Include CIB data")
	flag.BoolVar(&cliCollectorChecker, "collector.checker", false, "Include CheckerComponent data")
	flag.BoolVar(&cliCollectorInflux, "collector.influx", false, "Include InfluxDBWriter  data")
	flag.BoolVar(&cliCollectorInflux2, "collector.influx2", false, "Include InfluxDB2Writer data")
	flag.BoolVar(&cliCollectorGraphite, "collector.graphite", false, "Include GraphiteWriter data")
	flag.BoolVar(&cliCollectorOTLP, "collector.otlpmetrics", false, "Include OTLPMetricsWriter data")

	flag.BoolVar(&cliVersion, "version", false, "Print version")
	flag.BoolVar(&cliDebugLog, "debug", false, "Enable debug logging")

	flag.Parse()

	if cliVersion {
		fmt.Printf("icinga-exporter version: %s\n", buildVersion())
		os.Exit(0)
	}

	u, errURL := url.Parse(cliBaseURL)

	if errURL != nil {
		fmt.Fprintf(os.Stderr, "Invalid Icinga2 URL: %v", errURL)
	}

	logLevel := slog.LevelInfo

	if cliDebugLog {
		logLevel = slog.LevelDebug
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	}))

	if cliPassword == "" {
		if v, ok := os.LookupEnv("ICINGA2_EXPORTER_HTTP_PASSWORD"); ok {
			cliPassword = v
		}
	}

	// In general, listen to gosec. But it this case, I don't think someone
	// is going to overflow the uint TTL for the cache lifetime.
	// nolint:gosec
	cacheTTL := time.Duration(cliCacheTTL) * time.Second

	config := icinga.Config{
		BasicAuthUsername: cliUsername,
		BasicAuthPassword: cliPassword,
		CAFile:            cliCAFile,
		CertFile:          cliCertFile,
		KeyFile:           cliKeyFile,
		Insecure:          cliInsecure,
		CacheTTL:          cacheTTL,
		IcingaAPIURI:      *u,
	}

	c, errCli := icinga.NewClient(config)

	if errCli != nil {
		fmt.Fprintf(os.Stderr, "Could not create Icinga2 client : %v", errCli)
	}

	// Register Collectors
	prometheus.MustRegister(collector.NewIcinga2ApplicationCollector(c, logger))

	if cliCollectorCIB {
		prometheus.MustRegister(collector.NewIcinga2CIBCollector(c, logger))
	}

	if cliCollectorApiListener {
		prometheus.MustRegister(collector.NewIcinga2APICollector(c, logger))
	}

	if cliCollectorChecker {
		prometheus.MustRegister(collector.NewIcinga2CheckerCollector(c, logger))
	}

	if cliCollectorInflux {
		prometheus.MustRegister(collector.NewIcinga2InfluxDBCollector(c, logger))
	}

	if cliCollectorInflux2 {
		prometheus.MustRegister(collector.NewIcinga2InfluxDB2Collector(c, logger))
	}

	if cliCollectorGraphite {
		prometheus.MustRegister(collector.NewIcinga2GraphiteCollector(c, logger))
	}

	if cliCollectorOTLP {
		prometheus.MustRegister(collector.NewIcinga2OTLPMetricsCollector(c, logger))
	}

	// Create a central context to propagate a shutdown
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	srv := &http.Server{
		Addr:              cliListenAddress,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       3 * time.Second,
		WriteTimeout:      3 * time.Second,
		IdleTimeout:       5 * time.Second,
	}

	http.Handle(cliMetricsPath, promhttp.Handler())
	//nolint:errcheck
	http.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte(`
			<html>
			<head><title>Icinga2 Exporter</title></head>
			<body>
			<h1>Icinga2 Exporter</h1>
			<p><a href="` + cliMetricsPath + `">Metrics</a></p>
			</body>
			</html>`))
	})

	go func() {
		slog.Info("Listening on address", "port", cliListenAddress, "version", version, "commit", commit)
		// nolint:noinlineerr
		if err := srv.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			slog.Error("HTTP server error", "error", err.Error())
		}

		slog.Info("Received Shutdown. Stopped serving new connections.")
	}()

	// The signal channel will block until we registered signals are received.
	// We will then use a context with a timeout to shutdown the application gracefully.
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	shutdownCtx, shutdownCancel := context.WithTimeout(ctx, 5*time.Second)
	defer shutdownCancel()

	// We are using Shutdown() with a timeout to gracefully
	// shut down the server without interrupting any active connections.
	// nolint:noinlineerr
	if err := srv.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP shutdown error", "error", err.Error())
	}
}
