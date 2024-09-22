package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/Jacobbrewer1/satisfactory/pkg/logging"
	"github.com/google/subcommands"
)

func main() {
	if err := logging.GeneralLogger(appName); err != nil {
		fmt.Println("Error setting up logging:", err)
		os.Exit(1)
	}

	subcommands.Register(subcommands.HelpCommand(), "")
	subcommands.Register(subcommands.FlagsCommand(), "")
	subcommands.Register(subcommands.CommandsCommand(), "")

	subcommands.Register(new(versionCmd), "")
	subcommands.Register(new(startCmd), "")

	flag.Parse()

	// Listen for ctrl+c and kill signals
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, os.Interrupt, syscall.SIGTERM)
		got := <-sig
		slog.Info("Received signal, shutting down", slog.String("signal", got.String()))
		cancel()
	}()

	os.Exit(int(subcommands.Execute(ctx)))
}
