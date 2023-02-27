package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var issuesCmd = &cobra.Command{
	Use:   "issues [owner] [repository]",
	Short: "Queries the issues of a repository at owner/repository",
	Args:  cobra.ExactArgs(2),
	Run: func(cmd *cobra.Command, args []string) {
		if sinceFlag == "" {
			// Get since from the database and continue from there.
			since = metadata.IssuesUpdatedAt
		}

		// Query the issues.
		issues, err := gh.QueryIssues(owner, repository, since, limit)
		if err != nil {
			logger.Error("could not query issues", "error", err)
			os.Exit(1)
		}

		if format == "sql" {
			// Write the issues to the database.
			for _, i := range issues {
				_, err := db.AddIssue(i)
				if err != nil {
					logger.Error("could not add issue to database", "error", err)
					os.Exit(1)
				}

				for _, r := range i.Reactions {
					r.Issue = i.ID
					_, err := db.AddIssueReaction(r)
					if err != nil {
						logger.Error("could not add issue reaction to database", "error", err)
						os.Exit(1)
					}
				}

				for _, c := range i.Comments {
					c.Issue = i.ID
					_, err := db.AddIssueComment(c)
					if err != nil {
						logger.Error("could not add issue comment to database", "error", err)
						os.Exit(1)
					}

					for _, cr := range c.Reactions {
						cr.Issue = i.ID
						_, err := db.AddIssueCommentReaction(cr)
						if err != nil {
							logger.Error("could not add issue comment reaction to database", "error", err)
							os.Exit(1)
						}
					}
				}

				// Update the metadata.
				if i.UpdatedAt.After(metadata.IssuesUpdatedAt) {
					metadata.IssuesUpdatedAt = i.UpdatedAt
					metadata, err = db.AddMetadata(metadata)
					if err != nil {
						logger.Error("could not update metadata", "error", err)
						os.Exit(1)
					}
				}
			}
		} else {
			// Output the issues as JSON.
			outputJSON(issues, output)
		}
	},
}
