package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	cfgpkg "github.com/aalexanderkevin/getstarvio-backend/internal/config"
	dbplatform "github.com/aalexanderkevin/getstarvio-backend/internal/platform/db"
)

func newMigrateCommand() *cobra.Command {
	var rollbackOne bool

	cmd := &cobra.Command{
		Use:   "migrate",
		Short: "Run database migrations",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg := cfgpkg.MustLoad()
			db, err := dbplatform.Open(cfg)
			if err != nil {
				return fmt.Errorf("open db: %w", err)
			}
			defer func() {
				_ = dbplatform.Close(db)
			}()

			sqlDB, err := db.DB()
			if err != nil {
				return fmt.Errorf("db sql handle: %w", err)
			}

			if rollbackOne {
				return dbplatform.RollbackOne(sqlDB, cfg.DB.MigrationPath)
			}

			err = dbplatform.Migrate(sqlDB, cfg.DB.MigrationPath)
			if err != nil {
				fmt.Printf("Error when migration: %s \n", err.Error())
				os.Exit(1)
			}

			fmt.Println("Finish migrating database")
			return nil
		},
	}

	cmd.Flags().BoolVar(&rollbackOne, "down-one", false, "rollback one migration")
	return cmd
}
