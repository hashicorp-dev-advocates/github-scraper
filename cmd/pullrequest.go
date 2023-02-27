package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var pullrequestsCmd = &cobra.Command{
	Use:   "pullrequests [owner] [name]",
	Short: "Queries the pullrequests of a repository at owner/name",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if sinceFlag == "" {
			// Get since from the database and continue from there.
			since = metadata.PullrequestsUpdatedAt
		}

		// Query the pullrequests.
		pullrequests, err := gh.QueryPullrequests(owner, repository, since, limit)
		if err != nil {
			logger.Error("could not query pullrequests", "error", err)
			os.Exit(1)
		}

		if format == "sql" {
			// Write the pullrequests to the database.
			for _, p := range pullrequests {
				_, err := db.AddPullrequest(p)
				if err != nil {
					logger.Error("could not add pullrequest to database", "error", err)
					os.Exit(1)
				}

				for _, r := range p.Reactions {
					r.Pullrequest = p.ID
					_, err := db.AddPullrequestReaction(r)
					if err != nil {
						logger.Error("could not add pullrequest reaction to database", "error", err)
						os.Exit(1)
					}
				}

				for _, c := range p.Comments {
					c.Pullrequest = p.ID
					_, err := db.AddPullrequestComment(c)
					if err != nil {
						logger.Error("could not add pullrequest comment to database", "error", err)
						os.Exit(1)
					}

					for _, cr := range c.Reactions {
						cr.Pullrequest = p.ID
						_, err := db.AddPullrequestCommentReaction(cr)
						if err != nil {
							logger.Error("could not add pullrequest comment reaction to database", "error", err)
							os.Exit(1)
						}
					}
				}

				for _, r := range p.Reviews {
					r.Pullrequest = p.ID
					_, err := db.AddPullrequestReview(r)
					if err != nil {
						logger.Error("could not add pullrequest review to database", "error", err)
						os.Exit(1)
					}
				}

				for _, f := range p.Files {
					f.Pullrequest = p.ID
					_, err := db.AddPullrequestFile(f)
					if err != nil {
						logger.Error("could not add pullrequest file to database", "error", err)
						os.Exit(1)
					}
				}

				// Update the metadata.
				if p.UpdatedAt.After(metadata.PullrequestsUpdatedAt) {
					metadata.PullrequestsUpdatedAt = p.UpdatedAt
					metadata, err = db.AddMetadata(metadata)
					if err != nil {
						logger.Error("could not update metadata", "error", err)
						os.Exit(1)
					}
				}
			}
		} else {
			outputJSON(pullrequests, output)
		}
	},
}
