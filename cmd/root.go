package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/eveldcorp/devrel-github/config"
	"github.com/eveldcorp/devrel-github/database"
	"github.com/eveldcorp/devrel-github/github"
	"github.com/hashicorp/go-hclog"
	"github.com/spf13/cobra"
	"go.opentelemetry.io/otel/trace"
)

var cfg *config.Config

var db database.Database
var gh github.Github
var logger hclog.Logger
var tracer trace.Tracer
var err error

// CLI parameters.
var sinceFlag string
var since time.Time
var limit int
var format string
var output string

// CLI args.
var owner string
var repository string

var metadata database.Metadata

var rootCmd = &cobra.Command{
	Use:   "github",
	Short: "Queries GitHub for information.",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		// Config.
		cfg = config.New()

		// Logger.
		logger = hclog.Default()
		logger.SetLevel(hclog.Debug)

		owner = args[0]
		repository = args[1]

		// Output format.
		format, _ = cmd.Flags().GetString("format")
		if format != "sql" && format != "json" {
			logger.Error("unknown output format")
			os.Exit(1)
		}

		// Output location.
		output, _ = cmd.Flags().GetString("output")
		if format == "sql" && output == "" {
			logger.Error("no database connection string specified as output e.g. postgres://user:password@host:port/database")
			os.Exit(1)
		}

		// Start date for query.
		sinceFlag, _ = cmd.Flags().GetString("since")
		if sinceFlag != "" {
			since, err = time.Parse(time.RFC3339, sinceFlag)
			if err != nil {
				logger.Error("invalid time for since", "error", err)
				os.Exit(1)
			}
		}

		limit, _ = cmd.Flags().GetInt("limit")
		if limit < 0 {
			logger.Error("limit has to be a positive number")
			os.Exit(1)
		}

		// Database.
		if format == "sql" {
			db, err = database.New(output, cfg.DBmaxopenconns, cfg.DBconnmaxlifetime, logger)
			if err != nil {
				logger.Error("Could not connect to database", "error", err)
				os.Exit(1)
			}

			metadata, err = db.GetMetadata(owner, repository)
			if err != nil {
				logger.Error("could not query metadata", "error", err)
				os.Exit(1)
			}
		}

		// Github.
		gh = github.New(cfg.GitHubToken, logger)
	},
}

// Execute runs the main command.
func Execute() {
	// Parse parameters.
	rootCmd.PersistentFlags().StringVarP(&format, "format", "f", "json", "The format to output the data as, can be either json or sql")
	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "Where to write the data to")
	rootCmd.PersistentFlags().StringVarP(&sinceFlag, "since", "s", "", "Query data created after this date")
	rootCmd.PersistentFlags().IntVarP(&limit, "limit", "l", 0, "What to limit the number of results to")

	// Add subcommands.
	rootCmd.AddCommand(issuesCmd)
	rootCmd.AddCommand(pullrequestsCmd)
	rootCmd.AddCommand(releasesCmd)
	rootCmd.AddCommand(metricsCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Output the data as JSON.
func outputJSON(input interface{}, dest string) error {
	data, err := json.MarshalIndent(&input, "", " ")
	if err != nil {
		return fmt.Errorf("could not marshall json: %v+", err)
	}

	if dest == "" {
		fmt.Println(string(data))
	} else {
		err = os.WriteFile(dest, data, 0755)
		if err != nil {
			return fmt.Errorf("could not write json to file: %v+", err)
		}
	}

	return nil
}
