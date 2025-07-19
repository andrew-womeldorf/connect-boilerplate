package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/lmittmann/tint"
	slogformatter "github.com/samber/slog-formatter"
	"github.com/spf13/cobra"
)

func main() {
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var (
	verbose  bool
	jsonLogs bool
)

// RootCmd represents the base command when called without any subcommands
var RootCmd = &cobra.Command{
	Use:   "api",
	Short: "Example Connect RPC API CLI",
	Long: `Example Connect RPC API CLI provides commands for managing users.
It provides both server functionality and RPC client commands.`,
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Configure logger based on verbose flag
		logLevel := slog.LevelInfo
		if verbose {
			logLevel = slog.LevelDebug
		}

		// Configure logger
		var handler slog.Handler
		if jsonLogs {
			handler = slog.NewJSONHandler(os.Stderr, &slog.HandlerOptions{
				AddSource: true,
				Level:     logLevel,
			})
		} else {
			handler = tint.NewHandler(os.Stderr, &tint.Options{
				AddSource: true,
				Level:     logLevel,
			})
		}

		logger := slog.New(slogformatter.NewFormatterHandler(
			slogformatter.ErrorFormatter("error"),
		)(handler))
		slog.SetDefault(logger)
	},
	Run: func(cmd *cobra.Command, args []string) {
		// If no subcommand is provided, print help
		if err := cmd.Help(); err != nil {
			slog.Error("Failed to display help", "error", err)
			os.Exit(1)
		}
	},
}

func init() {
	// Add persistent flags that will be available to all commands
	RootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	RootCmd.PersistentFlags().BoolVar(&jsonLogs, "json", false, "Output logs in JSON format (default: text)")
}
