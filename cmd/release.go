package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var releasesCmd = &cobra.Command{
	Use:   "releases [owner] [name]",
	Short: "Queries the releases of a repository at owner/name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		// Query the releases.
		releases, err := gh.QueryReleases(owner, repository, limit)
		if err != nil {
			logger.Error("could not query releases", "error", err)
			os.Exit(1)
		}

		if format == "sql" {
			// Write the releases to the database.
			for _, r := range releases {
				_, err := db.AddRelease(r)
				if err != nil {
					logger.Error("could not add release to database", "error", err)
					os.Exit(1)
				}

				for _, a := range r.Assets {
					_, err := db.AddReleaseAsset(a)
					if err != nil {
						logger.Error("could not add release asset to database", "error", err)
						os.Exit(1)
					}
				}
			}
		} else {
			// Output the releases as JSON.
			outputJSON(releases, output)
		}
	},
}
