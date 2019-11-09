package main

import (
	"context"
	{{if .TLS}}"crypto/tls"
	"crypto/x509"
	"io/ioutil"{{end}}
	"flag"
	"fmt"
	{{if .Grpc}}"net"{{end}}
	{{if .HTTP}}"net/http"{{end}}
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	{{if .ClientHelper}}"{{ .ServicePath }}/client"{{end}}
	{{if .Jobs}}"{{ .ServicePath }}/jobs"{{end}}
	"{{ .ServicePath }}/pb/{{ package .ServiceName }}"
	"{{ .ServicePath }}/service"
	{{if .HTTP}}"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"{{end}}
	{{if .DB}}"github.com/go-pg/pg/v9"
	"github.com/go-pg/pg/v9/orm"{{end}}
	{{if .Jobs}}"github.com/gocraft/work"{{end}}
	{{if .RedisPool}}"github.com/gomodule/redigo/redis"{{end}}
	{{if .DB}}"github.com/google/uuid"{{end}}
	{{if .Nats}}"github.com/nats-io/nats.go"{{end}}
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	{{if .Grpc}}"google.golang.org/grpc"{{end}}
)

const serviceName = "{{ .ServiceName }}"

var (
	buildTag string

	appName = flag.String("app_name", "undefined", "global application name")

	{{if .HTTP}}httpPort    = flag.Int("http_port", 8888, "http port"){{end}}
	{{if .Grpc}}grpcPort    = flag.Int("grpc_port", 9200, "grpc port"){{end}}
	debug       = flag.Bool("debug", false, "enable debug mode")
	console     = flag.Bool("console_log", false, "change logs format from JSON to console output")
	{{if .TLS}}tlsRootCert = flag.String("tls_root_cert", "", "path to TLS root certificate"){{end}}
	{{if .DB}}
	dbHost     = flag.String("db_host", "postgres", "postgresql database name")
	dbPort     = flag.Int("db_port", 5432, "postgresql database port")
	dbName     = flag.String("db_name", "", "postgresql database name")
	dbUser     = flag.String("db_user", "", "postgresql database user name")
	dbPassword = flag.String("db_password", "", "postgresql database password")
	dbPoolSize = flag.Int("db_pool_size", 10, "postgresql connections pool size")
	{{end}}
	{{if .RedisPool}}redisHost = flag.String("redis_host", "redis://redis:6379", "redis host, e.g.: redis://redis:6379"){{end}}
	{{if .Nats}}
	natsHost         = flag.String("nats_host", "nats://nats-cluster:4222", "NATS host to connect, e.g.: nats://nats-cluster:4222")
	natsQueueSubject = flag.String("nats_queue_subject", serviceName, "NATS queue subject"){{end}}

	{{if .Grpc}}grpcServer *grpc.Server{{end}}
	{{if .HTTP}}httpServer *http.Server{{end}}
	{{if .Nats}}queueSubscr *nats.Subscription{{end}}
)

func main() {
	flag.Parse()

	// Set up logger
	logger := initLogger(*console, *debug, serviceName, buildTag)

	// Listen interrupt signal from OS
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(interrupt)

	// Context with the cancellation function
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Error group with custom context
	eg, ctx := errgroup.WithContext(ctx)

	{{if .TLS}}
	// TLS config for do services
	var tlsConfig *tls.Config
	if *tlsRootCert != "" {
		tlsConfig, err = loadTLSRootCert(*tlsRootCert)
		if err != nil {
			logger.Fatal().Err(err).Msg("tls root certificate loading")
		}
	}
	{{end}}

	{{if .DB}}
	// Set up datapase connection
	db := pg.Connect(&pg.Options{
		Addr:            fmt.Sprintf("%s:%d", *dbHost, *dbPort),
		User:            *dbUser,
		Password:        *dbPassword,
		Database:        *dbName,
		ApplicationName: *appName,
		PoolSize:        *dbPoolSize,
		TLSConfig:       tlsConfig,
	})
	defer db.Close()

	if *debug {
		// Log each database query
		db.AddQueryHook(newDBLogger(logger))
	}

	if err := db.CreateTable(&struct {
		tableName string    `pg:"test_tbl"`
		ID        uuid.UUID `pg:",pk"`
	}{}, &orm.CreateTableOptions{IfNotExists: true}); err != nil {
		logger.Fatal().Err(err).Msg("create table")
	}
	{{end}}

	{{if .RedisPool}}
	// Setup redis connections pool
	redisPool := newRedisPool(*redisHost, tlsConfig)
	{{end}}

	{{if .Jobs}}
	// Init worker pool
	pool := work.NewWorkerPool(jobs.WorkerPoolContext{}, 10, "{default}", redisPool)
	// Init worker pool context
	worker := jobs.NewWorker(pool)
	if err := worker.RegisterJobs(); err != nil {
		logger.Fatal().Err(err).Msg("worker jobs registration")
	}
	// Start processing jobs
	pool.Start()
	// Stop the pool
	defer pool.Stop()
	// Init enqueuer
	enqueuer := work.NewEnqueuer("{default}", redisPool)
	{{end}}

	{{if .GrpcClients}}{{range $value := .GrpcClients}}
	// Set up a connection to the server.
	{{urlTovar $value}}Conn, err := grpc.Dial("{{$value}}", grpc.WithInsecure())
	if err != nil {
		logger.Fatal().Err(err).Str("host", "{{$value}}").Msg("grpc client connection error")
	}
	{{urlTovar $value}}Client := {{ package $value }}.NewServiceClient({{urlTovar $value}}Conn)
	{{end}}{{end}}

	{{if .Nats}}
	// Set up NATS connection
	natsOpts := setupNatsConnOptions(logger, []nats.Option{nats.Name(serviceName)})
	nc, err := nats.Connect(*natsHost, natsOpts...)
	if err != nil {
		logger.Fatal().Err(err).Str("host", *natsHost).Msg("nats client connection error")
	}
	defer nc.Close()
	{{end}}

	{{if .Sub}}
	// Subscribe to {{.ServiceName}}.* channel in {{.ServiceName}} queue
	queueSubscr, err := nc.QueueSubscribe(*natsQueueSubject, serviceName, func(msg *nats.Msg) {
		logger.Info().
			Str("subject", msg.Subject).
			Str("reply", msg.Reply).
			Str("data", string(msg.Data))).
			Msg("received message")
	})
	if err != nil {
		logger.Fatal().Err(err).Msg("nats queue subscription error")
	}
	defer queueSubscr.Unsubscribe()
	{{end}}

	{{if .Pub}}
	eg.Go(func() error {
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()
		for {
			select {
			case {{ asis "<-" }}ticker.C:
				msg := fmt.Sprintf("publisher {{ .ServiceName }}: Current time is %s", time.Now().String())
				logger.Debug().Msg(msg)
				nc.Publish(*natsQueueSubject, []byte(msg))
			case {{ asis "<-" }}ctx.Done():
				return nil
			}
		}
	})
	{{end}}

	{{if .RPC}}
	// Init RPC service
	{{ camel .ServiceName }}Service := service.New(logger{{if .DB}}, db{{end}}{{if .Jobs}}, enqueuer{{end}})
	{{end}}

	{{if .Grpc}}
	// Set up grpc server
	grpcServer = grpc.NewServer()
	// Registering of grpc services
	{{ package .ServiceName }}.RegisterServiceServer(grpcServer, {{ camel .ServiceName }}Service)
	// Run gRPC server
	eg.Go(func() error {
		addr := fmt.Sprintf(":%d", *grpcPort)
		grpcLis, err := net.Listen("tcp", addr)
		if err != nil {
			return errors.Wrap(err, "grpc failed to listen")
		}
		logger.Info().Str("listen", addr).Msg("gRPC server started")
		if err := grpcServer.Serve(grpcLis); err != nil {
			return errors.Wrap(err, "grpc server")
		}
		return nil
	})
	{{end}}

	{{if .HTTP}}
	router := chi.NewRouter()
	router.Use(middleware.Recoverer)
	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNoContent)
	})

	httpServer = &http.Server{
		Addr:         fmt.Sprintf(":%d", *httpPort),
		Handler:      router,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}
	{{end}}

	{{if .Twirp}}
	// Registering of twirp rpc services
	{{ camel .ServiceName }}Handler := {{ package .ServiceName }}.NewServiceServer({{ camel .ServiceName }}Service, nil)
	router.Mount({{ camel .ServiceName }}Handler.PathPrefix(), {{ camel .ServiceName }}Handler)
	{{end}}

	{{if .HTTP}}
	// Run http server
	eg.Go(func() error {
		logger.Info().Str("listen", httpServer.Addr).Msg("HTTP server started")
		if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
			return errors.Wrap(err, "http server")
		}
		return nil
	})
	{{end}}

	// Wait for interrupt signal or context cancellation
	select {
	case {{ asis "<-" }}interrupt:
		break
	case {{ asis "<-" }}ctx.Done():
		break
	}

	logger.Info().Msg("received shutdown signal")

	{{if .Grpc}}
	if grpcServer != nil {
		grpcServer.GracefulStop()
	}
	{{end}}

	{{if .HTTP}}
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()
	if httpServer != nil {
		_ = httpServer.Shutdown(shutdownCtx)
	}
	{{end}}

	if err := eg.Wait(); err != nil {
		logger.Fatal().Err(err).Msg("failed to wait goroutine group")
	}

	logger.Info().Msg("service stopped")
}

{{if .TLS}}
func loadTLSRootCert(path string) (*tls.Config, error) {
	caPool := x509.NewCertPool()
	severCert, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}
	caPool.AppendCertsFromPEM(severCert)
	return &tls.Config{
		RootCAs:            caPool,
		InsecureSkipVerify: true,
	}, nil
}
{{end}}
func initLogger(console, debug bool, serviceName, buildTag string) (logger zerolog.Logger) {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}

	if console {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatLevel = func(i interface{}) string {
			return strings.ToUpper(fmt.Sprintf("| %-6s|", i))
		}
		logger = zerolog.New(output).With().Timestamp().Logger()
	} else {
		logger = log.Logger
	}

	logger = logger.With().Str("service", serviceName).Str("build", buildTag).Logger()
	logger.Debug().Msg("debug mode enabled")

	return logger
}