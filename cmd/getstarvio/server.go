package main

import (
	"fmt"

	"github.com/spf13/cobra"

	_ "github.com/aalexanderkevin/getstarvio-backend/docs"
	"github.com/aalexanderkevin/getstarvio-backend/internal/app"
	cfgpkg "github.com/aalexanderkevin/getstarvio-backend/internal/config"
	httpserver "github.com/aalexanderkevin/getstarvio-backend/internal/http"
	dbplatform "github.com/aalexanderkevin/getstarvio-backend/internal/platform/db"
)

func newServerCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Run HTTP API server",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cfgpkg.MustLoad()
			db, err := dbplatform.Open(cfg)
			if err != nil {
				return fmt.Errorf("open db: %w", err)
			}
			defer func() {
				_ = dbplatform.Close(db)
			}()

			container := app.NewContainer(cfg, db)
			srv := httpserver.NewServer(container)
			return srv.Start()
		},
	}
}
