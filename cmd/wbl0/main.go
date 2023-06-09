package main

import (
	"context"
	"fmt"

	"github.com/AlexeyBazhin/wbL0/internal/api/server"
	"github.com/AlexeyBazhin/wbL0/internal/api/stanListener"
	"github.com/AlexeyBazhin/wbL0/internal/cache"
	"github.com/AlexeyBazhin/wbL0/internal/db"
	"github.com/AlexeyBazhin/wbL0/internal/domain/service"
	"github.com/AlexeyBazhin/wbL0/internal/repository"
	"github.com/nats-io/stan.go"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"github.com/kelseyhightower/envconfig"
	_ "github.com/lib/pq"
)

// type config struct {
// 	Env             string `required:"true" default:"development" desc:"production, development"`
// 	DSN             string `required:"true" default:":5432" desc:"DSN для соединения с базой данных"`
// 	BindAddr        string `required:"true" default:":8080" split_words:"true" desc:"Адрес и порт входящих соединений"`
// 	ReadTimeout     int    `required:"true" default:"10" split_words:"true" desc:"Таймаут на чтение запроса"`
// 	WriteTimeout    int    `required:"true" default:"10" split_words:"true" desc:"Таймаут на запись ответа"`
// 	ShutdownTimeout int    `required:"true" default:"30" split_words:"true" desc:"Время до принудительного завершения сервиса после получения сигнала выхода (s)"`
// }

func main() {
	var (
		ctx, ctxCancel = context.WithCancel(context.Background())
		cfg            = new(config)
		logger         *zap.Logger
		err            error
	)
	defer ctxCancel()

	if logger, err = zap.NewDevelopment(); err != nil {
		panic(fmt.Errorf("[main] Unable to start logger: %w", err))
	}
	sugaredLogger := logger.Sugar()

	if err := envconfig.Process("APP", cfg); err != nil {
		sugaredLogger.Fatal(err.Error())
	}

	connectedDB, err := db.ConnectPostgreSQL(ctx, cfg.DSN)
	if err != nil {
		panic(fmt.Errorf("[main] Unable to open connection with DB: %w", err))
	}
	sugaredLogger.Info("[main] Successful connection to DB")

	redisClient := redis.NewClient(&redis.Options{
		Addr:     "cache:6379",
		Password: "",
		DB:       0,
	})
	if pong, err := redisClient.Ping(ctx).Result(); err != nil {
		panic(fmt.Errorf("[main] Unable to open connection with Redis: %w", err))
	} else {
		sugaredLogger.Infof("[main] Successful connection to Redis %v", pong)
	}

	repo := repository.NewRepository(connectedDB, sugaredLogger)
	cache := cache.NewCache(redisClient)
	svc := service.NewService(repo, cache)

	errGroup, egCtx := errgroup.WithContext(ctx)
	errGroup.Go(
		server.New(
			egCtx, //server.WithContext(egCtx),
			server.WithService(svc),
			server.WithBindAddress(cfg.BindAddr),
			server.WithLogger(sugaredLogger),
		).Run())

	stanConn, err := stan.Connect("amethyst-cluster", "wbL0", stan.NatsURL("http://nats:4222"))
	if err != nil {
		panic(fmt.Errorf("[main] Unable to open connection with message broker nats-streaming: %w", err))
	}
	sugaredLogger.Info("[main] Successful connection to nats-streaming ")
	defer stanConn.Close()

	errGroup.Go(
		stanListener.New(
			stanListener.WithService(svc),
			stanListener.WithStanConn(stanConn),
			stanListener.WithSubject("models"),
			stanListener.WithContext(egCtx),
			stanListener.WithLogger(sugaredLogger),
		).Run())

	// errGroup.Go(
	// 	func() error {
	// 		shutdown := make(chan os.Signal, 1)
	// 		signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	// 		sugaredLogger.Info("[signal-watcher] started")

	// 		select {
	// 		case sig := <-shutdown:
	// 			return fmt.Errorf("terminated with signal: %s", sig.String())
	// 		case <-ctx.Done():
	// 			return nil
	// 		}
	// 	})

	if err := errGroup.Wait(); err != nil {
		sugaredLogger.Error("successful shutdown", err)
	}
}
