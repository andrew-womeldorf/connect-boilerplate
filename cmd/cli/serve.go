package main

import (
	"log/slog"
	"os"

	"github.com/spf13/cobra"

	"github.com/andrew-womeldorf/connect-boilerplate/internal/logging"
	"github.com/andrew-womeldorf/connect-boilerplate/internal/server"
)

var (
	port       int
	jsonFormat bool
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Start the API server",
	Long:  `Start the API server that provides Connect RPC endpoints.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Setup logger
		logging.SetupLogger(jsonFormat, slog.LevelInfo)

		// Create and run server
		srv := server.NewServer(port)
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
	serveCmd.Flags().BoolVar(&jsonFormat, "json", false, "Use JSON log format")
}
