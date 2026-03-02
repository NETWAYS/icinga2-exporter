package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"net/url"
	"os"
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

	flag.BoolVar(&cliVersion, "version", false, "Print version")
	flag.BoolVar(&cliDebugLog, "debug", false, "Enable debug logging")

	flag.Parse()

	if cliVersion {
		fmt.Printf("icinga-exporter version: %s\n", version)
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
	prometheus.MustRegister(collector.NewIcinga2CIBCollector(c, logger))
	prometheus.MustRegister(collector.NewIcinga2ApplicationCollector(c, logger))

	if cliCollectorApiListener {
		prometheus.MustRegister(collector.NewIcinga2APICollector(c, logger))
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

	log.Printf("Version: %s", buildVersion())
	log.Printf("Listening on address: %s", cliListenAddress)
	log.Fatal(http.ListenAndServe(cliListenAddress, nil))
}
