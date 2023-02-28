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

	"lms/students_svc/client"
	db "lms/students_svc/db/sqlc"
	"lms/students_svc/service"
	"lms/students_svc/utils"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/oklog/oklog/pkg/group"
	stdprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
	"github.com/go-kit/log"
	_ "github.com/lib/pq"
)

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		fmt.Println(err)
		// log.Fatal().Err(err).Msg("cannot load config")
	}

	fs := flag.NewFlagSet("students_svc", flag.ExitOnError)
	var (
		debugAddr = fs.String("debug.addr", ":8080", "Debug and metrics listen address")
		httpAddr  = fs.String("http-addr", ":8081", "HTTP listen address")
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
			Subsystem: "students_service",
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, []string{"method"})
	}
	var duration metrics.Histogram
	{
		// Endpoint-level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: "api",
			Subsystem: "students_service",
			Name:      "request_duration_seconds",
			Help:      "Request duration in seconds.",
		}, []string{"method"})
	}
	http.DefaultServeMux.Handle("/metrics", promhttp.Handler())

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		logger.Log("cannot connect to DB", err)
	}

	runDBMigration(config.MigrationURL, config.DBSource, logger)

	store := db.New(conn)

	courseSvc, err := client.NewHTTPClient(config.CoursesHTTPServerAddress, logger)
	if err != nil {
		logger.Log("transport", "debug/HTTP", "during", "Client", "err", "cannot connect to course client", err)
	}

	var (
		std_service = service.New(store, courseSvc, logger, count, duration)
		endpoints   = service.MakeServerEndpoints(std_service, logger, duration)
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

func runDBMigration(migrationURL string, dbSource string, logger log.Logger) {
	migration, err := migrate.New(migrationURL, dbSource)
	if err != nil {
		logger.Log("cannot create new migrate instance", err)
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		logger.Log("failed to run migrate up", err)
	}

	logger.Log("db migrated successfully")
}
