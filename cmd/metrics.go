package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var metricsCmd = &cobra.Command{
	Use:   "metrics [owner] [name]",
	Short: "Queries the metrics of a repository at owner/name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Query the metrics.
		metrics, err := gh.QueryMetrics(owner, repository)
		if err != nil {
			logger.Error("could not query metrics", "error", err)
			os.Exit(1)
		}

		if format == "sql" {
			// Write the metrics to the database.
			db.AddMetrics(metrics)

			db.AddTrafficClones(metrics.Clones)
			db.AddTrafficViews(metrics.Views)

			for _, r := range metrics.Referrers {
				db.AddTrafficReferrer(r)
			}

			for _, p := range metrics.Paths {
				db.AddTrafficPath(p)
			}
		} else {
			// Output the releases as JSON.
			outputJSON(metrics, output)
		}
	},
}
