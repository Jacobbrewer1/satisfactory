package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"runtime"

	"github.com/Jacobbrewer1/goredis/redis"
	"github.com/Jacobbrewer1/satisfactory/pkg/alerts"
	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	svc "github.com/Jacobbrewer1/satisfactory/pkg/services/watcher"
	uhttp "github.com/Jacobbrewer1/satisfactory/pkg/utils/http"
	"github.com/Jacobbrewer1/vaulty/pkg/vaulty"
	"github.com/google/subcommands"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/spf13/viper"
)

type startCmd struct {
	// port is the port to listen on
	port string

	// configLocation is the location of the config file
	configLocation string
}

func (s *startCmd) Name() string {
	return "start"
}

func (s *startCmd) Synopsis() string {
	return "Start the bot"
}

func (s *startCmd) Usage() string {
	return `start:
  Start the bot.
`
}

func (s *startCmd) SetFlags(f *flag.FlagSet) {
	f.StringVar(&s.port, "port", "8080", "The port to listen on")
	f.StringVar(&s.configLocation, "config", "config.json", "The location of the config file")
}

func (s *startCmd) Execute(ctx context.Context, f *flag.FlagSet, _ ...interface{}) subcommands.ExitStatus {
	r := mux.NewRouter()
	service, err := s.setup(ctx, r)
	if err != nil {
		slog.Error("Error setting up bot", slog.String(logging.KeyError, err.Error()))
		return subcommands.ExitFailure
	}

	slog.Info(
		"Starting application",
		slog.String("version", Commit),
		slog.String("runtime", fmt.Sprintf("%s %s/%s", runtime.Version(), runtime.GOOS, runtime.GOARCH)),
		slog.String("build_date", Date),
	)

	srv := &http.Server{
		Addr:    ":" + s.port,
		Handler: r,
	}

	go func(service svc.Service) {
		err := service.Start()
		if err != nil {
			slog.Error("Error starting bot", slog.String(logging.KeyError, err.Error()))
			os.Exit(1)
		}
	}(service)

	// Start the server in a goroutine, so we can listen for the context to be done.
	go func(srv *http.Server) {
		err := srv.ListenAndServe()
		if errors.Is(err, http.ErrServerClosed) {
			slog.Info("Server closed gracefully")
			os.Exit(0)
		} else if err != nil {
			slog.Error("Error serving requests", slog.String(logging.KeyError, err.Error()))
			os.Exit(1)
		}
	}(srv)

	<-ctx.Done()
	slog.Info("Shutting down application")
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("Error shutting down application", slog.String(logging.KeyError, err.Error()))
		return subcommands.ExitFailure
	}

	return subcommands.ExitSuccess
}

func (s *startCmd) setup(ctx context.Context, r *mux.Router) (service svc.Service, err error) {
	v := viper.New()
	v.SetConfigFile(s.configLocation)
	if err := v.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("error reading config file: %w", err)
	}

	if !v.IsSet("vault") {
		return nil, errors.New("vault configuration not found")
	}

	slog.Info("Vault configuration found, attempting to connect")
	vc, err := vaulty.NewClient(
		vaulty.WithContext(ctx),
		vaulty.WithGeneratedVaultClient(v.GetString("vault.address")),
		vaulty.WithUserPassAuth(
			v.GetString("vault.auth.username"),
			v.GetString("vault.auth.password"),
		),
		vaulty.WithKvv2Mount(v.GetString("vault.kvv2_mount")),
	)
	if err != nil {
		return nil, fmt.Errorf("error creating vault client: %w", err)
	}

	slog.Debug("Vault client created")

	vs, err := vc.GetKvSecretV2(ctx, v.GetString("vault.bot.secret_name"))
	if errors.Is(err, vaulty.ErrSecretNotFound) {
		return nil, fmt.Errorf("secrets not found in vault: %s", v.GetString("vault.bot.token_path"))
	} else if err != nil {
		return nil, fmt.Errorf("error getting secrets from vault: %w", err)
	}

	if err := redis.NewPool(
		redis.WithDefaultPool(),
		redis.FromViper(v)...,
	); err != nil {
		return nil, fmt.Errorf("error creating redis pool: %w", err)
	}

	am := alerts.NewDiscordManager(vs.Data[v.GetString("vault.bot.alerts_url_key")].(string))
	service = svc.NewService(ctx, am, v.GetString("redis.info_list_name"), v.GetString("redis.details_list_name"))

	r.HandleFunc("/metrics", uhttp.InternalOnly(promhttp.Handler())).Methods(http.MethodGet)
	r.HandleFunc("/health", uhttp.InternalOnly(healthHandler())).Methods(http.MethodGet)

	return service, nil
}
