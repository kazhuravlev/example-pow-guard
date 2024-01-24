package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/kazhuravlev/example-pow-guard/internal/api"
	"github.com/kazhuravlev/example-pow-guard/internal/facade"
	"github.com/kazhuravlev/example-pow-guard/pkg/randstring"
	"log/slog"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	if err := cmdRun(); err != nil {
		fmt.Println(strings.Repeat("=", 80))
		fmt.Printf("The sky is falling: %v\n", err)
		fmt.Println(strings.Repeat("=", 80))
	}
}

func cmdRun() error {
	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT)

	defer cancel()

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		AddSource:   false,
		Level:       slog.LevelDebug,
		ReplaceAttr: nil,
	}))

	var flagPort int
	flag.IntVar(&flagPort, "port", 8888, "listen port")
	if err := flag.CommandLine.Parse(os.Args[1:]); err != nil {
		return fmt.Errorf("parse command flags: %w", err)
	}

	pRng := rand.New(rand.NewSource(time.Now().Unix()))
	rndQuotes, err := randstring.New(pRng, quotes)
	if err != nil {
		return fmt.Errorf("init random strings store: %w", err)
	}

	apiInst, err := api.New(logger, rndQuotes)
	if err != nil {
		return fmt.Errorf("init api instance: %w", err)
	}

	facadeInst, err := facade.New(logger, apiInst, flagPort)
	if err != nil {
		return fmt.Errorf("init facade instance: %w", err)
	}

	if err := facadeInst.Run(ctx); err != nil {
		return fmt.Errorf("run facade instance: %w", err)
	}

	select {
	case <-ctx.Done():
		logger.Warn("server is going to shutdown")
		logger.Info("wait all connections to stop")
		facadeInst.Wait()
		logger.Info("all connections was stopped")
	}

	return nil
}
