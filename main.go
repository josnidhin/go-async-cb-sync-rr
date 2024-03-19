/**
 * @author Jose Nidhin
 */
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/redis/go-redis/v9"
)

const (
	ChannelName = "async-sync-channel"
)

var logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

var httpServer *http.Server
var wg sync.WaitGroup

type AppFlag struct {
	Port      uint
	RedisAddr string
}

type RedisClient interface {
	redis.Cmdable

	Subscribe(ctx context.Context, channel ...string) *redis.PubSub
}

func main() {
	slog.SetDefault(logger)
	ctx, cancel := context.WithCancel(context.Background())

	appFlags := setupAppFlag()

	redisClient, err := initRedis(appFlags.RedisAddr)
	if err != nil {
		panic(err)
	}

	router := chi.NewRouter()
	router.Use(RequestId())
	router.Use(RequestLogger(logger))
	router.Use(middleware.AllowContentType("application/json"))
	router.NotFound(func(w http.ResponseWriter, r *http.Request) {
		logger.Info("route not found")
		w.WriteHeader(http.StatusNotFound)
	})
	router.Post("/sync", SyncHandler(redisClient, logger))
	router.Route("/callback", func(r chi.Router) {
		r.Post("/{id}", CallbackHandler(redisClient, logger))
	})

	httpServer = &http.Server{
		Addr:    ":" + strconv.FormatUint(uint64(appFlags.Port), 10),
		Handler: router,
	}

	wg.Add(1)
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer wg.Done()
		logger.Info("HTTP Server started", slog.String("address", httpServer.Addr))
		err := httpServer.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			logger.Error("Unexpected HTTP server error", slog.Any("error", err))
			shutdown(ctx, cancel)
		}
	}(ctx, cancel)

	wg.Add(1)
	go func(ctx context.Context, cancel context.CancelFunc) {
		defer wg.Done()
		signalHandler(ctx, cancel)
	}(ctx, cancel)

	wg.Wait()
}

func setupAppFlag() *AppFlag {
	var port uint
	var redisAddr string

	flag.UintVar(&port, "port", 8080, "HTTP port for the service. eg: 8080 (default)")

	flag.StringVar(&redisAddr, "msgtype", "localhost:6379", "Redis DB address, eg: localhost:6379 (default)")

	flag.Parse()

	return &AppFlag{
		Port:      port,
		RedisAddr: redisAddr,
	}
}

func initRedis(redisAddr string) (client RedisClient, err error) {
	opts := &redis.Options{
		Addr: redisAddr,
	}

	client = redis.NewClient(opts)
	statusCmd := client.Ping(context.Background())
	err = statusCmd.Err()
	if err != nil {
		err = fmt.Errorf("redis connection failed: %w", err)
	}
	logger.Info("Redis connected")
	return
}

func signalHandler(ctx context.Context, cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case s := <-sigChan:
		logger.Info("Shutdown signal received",
			slog.String("signal", s.String()))
		break
	case <-ctx.Done():
		break
	}

	shutdown(ctx, cancel)
}

func shutdown(ctx context.Context, cancel context.CancelFunc) {
	defer cancel()
	logger.Info("Graceful shutdown initialized")

	var err error

	err = httpServer.Shutdown(ctx)
	if err != nil {
		logger.Error("HTTP Server shutdown error", slog.Any("error", err))
	}
}
