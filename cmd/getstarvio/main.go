// @title Getstarvio Backend API
// @version 1.0
// @description Getstarvio backend API documentation.
// @BasePath /
// @schemes http https
//
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Bearer access token, e.g. "Bearer <token>"
package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

func newRootCommand() *cobra.Command {
	root := &cobra.Command{
		Use:   "getstarvio",
		Short: "Getstarvio backend service",
	}

	root.AddCommand(newServerCommand())
	root.AddCommand(newMigrateCommand())
	root.AddCommand(newWorkerCommand())

	return root
}

func main() {
	if err := newRootCommand().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
