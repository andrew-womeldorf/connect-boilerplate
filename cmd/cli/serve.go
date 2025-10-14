package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/andrew-womeldorf/connect-boilerplate/internal/server"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/services/user/store/sqlite"
)

var port int

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long:  `Start the API server that provides Connect RPC endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		store, err := sqlite.NewStore(context.Background(), ":memory:")
		if err != nil {
			panic(err)
		}

		// Create and run server
		srv := server.NewServer(port, store)
		if err := srv.Run(); err != nil {
			slog.Error("Failed to run server", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	RootCmd.AddCommand(serveCmd)

	// Add flags specific to the serve command
	serveCmd.Flags().IntVarP(&port, "port", "p", 8088, "Port to listen on")
}
