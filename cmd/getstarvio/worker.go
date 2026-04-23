package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	cfgpkg "github.com/aalexanderkevin/getstarvio-backend/internal/config"
	"github.com/aalexanderkevin/getstarvio-backend/internal/modules/reminder"
	dbplatform "github.com/aalexanderkevin/getstarvio-backend/internal/platform/db"
	"github.com/aalexanderkevin/getstarvio-backend/internal/platform/meta"
)

func newWorkerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "worker",
		Short: "Run reminder worker",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cfgpkg.MustLoad()
			db, err := dbplatform.Open(cfg)
			if err != nil {
				return fmt.Errorf("open db: %w", err)
			}
			defer func() {
				_ = dbplatform.Close(db)
			}()

			repo := reminder.NewRepo(db)
			svc := reminder.NewService(repo, meta.NewClient(cfg.Meta, meta.NewGormFacebookLogStore(db)), cfg.Meta)

			ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
			defer cancel()

			poll := time.Duration(cfg.Worker.PollIntervalSeconds) * time.Second
			if poll < 5*time.Second {
				poll = 5 * time.Second
			}

			return svc.RunWorker(ctx, poll)
		},
	}
}
