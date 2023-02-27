package main

import (
	"database/sql"
	"flag"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"courses/client"
	db "courses/db/sqlc"
	"courses/service"
	"courses/utils"

	"github.com/oklog/oklog/pkg/group"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	_ "github.com/lib/pq"
)

func main() {
	fs := flag.NewFlagSet("courses_svc", flag.ExitOnError)
	var (
		debugAddr = fs.String("debug.addr", ":7070", "Debug and metrics listen address")
		httpAddr  = fs.String("http-addr", ":7071", "HTTP listen address")
	)
	flag.Parse()

	var logger log.Logger
	{
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	var count metrics.Counter
	{
		count = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: "api",
			Subsystem: "courses_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"})
	}
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "courses_service",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	config, err := utils.LoadConfig(".")
	if err != nil {
		logger.Log("config", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Log("db conn error", err)
	}

	store := db.New(conn)

	studentsSvc, err := client.NewHTTPClient(config.StudentsHTTPServerAddress, logger)
	if err != nil {
		logger.Log("cannot connect to students client", err)
	}

	var (
		crs_service = service.New(store, studentsSvc, logger, count, duration)
		endpoints   = service.MakeServerEndpoints(crs_service, logger, duration)
		httpHandler = service.MakeHTTPHandler(endpoints, logger)
	)

	var g group.Group
	{
		debugListener, err := net.Listen("tcp", *debugAddr)
		if err != nil {
			logger.Log("transport", "debug/HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "debug/HTTP", "addr", *debugAddr)
			return http.Serve(debugListener, http.DefaultServeMux)
		}, func(error) {
			debugListener.Close()
		})
	}
	{
		// The HTTP listener mounts the Go kit HTTP handler we created.
		httpListener, err := net.Listen("tcp", *httpAddr)
		if err != nil {
			logger.Log("transport", "HTTP", "during", "Listen", "err", err)
			os.Exit(1)
		}
		g.Add(func() error {
			logger.Log("transport", "HTTP", "addr", *httpAddr)
			return http.Serve(httpListener, httpHandler)
		}, func(error) {
			httpListener.Close()
		})
	}
	{
		// This function just sits and waits for ctrl-C.
		cancelInterrupt := make(chan struct{})
		g.Add(func() error {
			c := make(chan os.Signal, 1)
			signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
			select {
			case sig := <-c:
				return fmt.Errorf("received signal %s", sig)
			case <-cancelInterrupt:
				return nil
			}
		}, func(error) {
			close(cancelInterrupt)
		})
	}
	logger.Log("exit", g.Run())
}
