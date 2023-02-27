package github

import (
	"context"
	"fmt"
	"time"

	"github.com/eveldcorp/devrel-github/database"
	githubv3 "github.com/google/go-github/v50/github"
	"github.com/hashicorp/go-hclog"
	"github.com/hashicorp/go-retryablehttp"
	"github.com/shurcooL/githubv4"
	"golang.org/x/oauth2"
)

type Github interface {
	QueryIssues(owner string, repository string, since time.Time, limit int) ([]database.Issue, error)
	QueryPullrequests(owner string, repository string, since time.Time, limit int) ([]database.Pullrequest, error)
	QueryReleases(owner string, repository string, limit int) ([]database.Release, error)
	QueryMetrics(owner string, repository string) (database.Metrics, error)
}

type GithubImpl struct {
	v4     *githubv4.Client
	v3     *githubv3.Client
	logger hclog.Logger
}

func New(token string, logger hclog.Logger) Github {
	logger = logger.Named("Github")

	src := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: token})

	retryClient := retryablehttp.NewClient()
	retryClient.RetryMax = 10

	httpClient := retryClient.StandardClient()
	httpClient.Transport = &oauth2.Transport{
		Base:   httpClient.Transport,
		Source: oauth2.ReuseTokenSource(nil, src),
	}

	v4 := githubv4.NewClient(httpClient)
	v3 := githubv3.NewClient(httpClient)

	return &GithubImpl{
		v4:     v4,
		v3:     v3,
		logger: logger,
	}
}

func (g *GithubImpl) QueryIssues(owner string, repository string, since time.Time, limit int) ([]database.Issue, error) {
	var query struct {
		Repository struct {
			Issues struct {
				Nodes      []GithubIssue
				PageInfo   PageInfo
				TotalCount int
			} `graphql:"issues(filterBy: {since: $since}, first: 50, after: $cursor, orderBy: { field:UPDATED_AT, direction: ASC })"`
		} `graphql:"repository(name: $repository, owner: $owner)"`
		RateLimit RateLimit
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repository": githubv4.String(repository),
		"since":      githubv4.DateTime{Time: since},
		"cursor":     (*githubv4.String)(nil),
	}

	var ratelimit RateLimit

	page := 0
	issues := []database.Issue{}

	for {
		g.logger.Debug("Querying issues", "owner", owner, "repository", repository, "page", page)
		err := g.v4.Query(context.Background(), &query, variables)
		if err != nil {
			return issues, err
		}

		if ratelimit.Cost > ratelimit.Remaining {
			return issues, fmt.Errorf("the query would exceed the current rate limit")
		}

		// Process issues
		for _, i := range query.Repository.Issues.Nodes {
			issue := database.Issue{
				ID:                i.ID,
				Owner:             owner,
				Repository:        repository,
				Author:            i.Author.Login,
				AuthorAssociation: i.AuthorAssociation,
				Title:             i.Title,
				Body:              i.Body,
				CreatedAt:         i.CreatedAt,
				PublishedAt:       i.PublishedAt,
				UpdatedAt:         i.UpdatedAt,
				LastEditedAt:      i.LastEditedAt,
				ClosedAt:          i.ClosedAt,
				State:             i.State,
				Locked:            i.Locked,
				Closed:            i.Closed,
				Reactions:         []database.IssueReaction{},
				Comments:          []database.IssueComment{},
				Labels:            database.StringArray{},
			}

			timelineItems := i.TimelineItems.Nodes
			if len(timelineItems) > 0 {
				// Get the last actor that closed the issue
				issue.ClosedBy = timelineItems[len(timelineItems)-1].IssueTimelineItemsConnection.Actor.Login
			}

			// Reactions
			for _, r := range i.ReactionGroups {
				reaction := database.IssueReaction{
					Reaction: r.Content,
					Count:    r.Users.TotalCount,
				}
				issue.Reactions = append(issue.Reactions, reaction)
			}

			// Labels
			for _, l := range i.Labels.Nodes {
				issue.Labels = append(issue.Labels, l.Name)
			}

			// Query additional comments
			if i.Comments.PageInfo.HasNextPage {
				g.logger.Debug("Need to query additional comments")
				comments, err := g.QueryIssueComments(owner, repository, i.Number, i.Comments.PageInfo.EndCursor)
				if err != nil {
					return issues, err
				}
				i.Comments.Nodes = append(i.Comments.Nodes, comments...)
			}
			// Comments
			for _, c := range i.Comments.Nodes {
				comment := database.IssueComment{
					ID:                c.ID,
					Author:            c.Author.Login,
					AuthorAssociation: c.AuthorAssociation,
					Body:              c.Body,
					CreatedAt:         c.CreatedAt,
					LastEditedAt:      c.LastEditedAt,
					PublishedAt:       c.PublishedAt,
					UpdatedAt:         c.UpdatedAt,
				}

				// Comment reactions
				for _, cr := range c.ReactionGroups {
					reaction := database.IssueCommentReaction{
						Content: cr.Content,
						Count:   cr.Users.TotalCount,
					}
					comment.Reactions = append(comment.Reactions, reaction)
				}
				issue.Comments = append(issue.Comments, comment)
			}

			issues = append(issues, issue)

			if len(issues) == limit {
				return issues, nil
			}
		}

		if !query.Repository.Issues.PageInfo.HasNextPage {
			break
		}

		ratelimit = query.RateLimit

		variables["cursor"] = githubv4.String(query.Repository.Issues.PageInfo.EndCursor)
		page++
	}

	return issues, nil
}

func (g *GithubImpl) QueryIssueComments(owner string, repository string, number int, cursor string) ([]GithubComment, error) {
	var query struct {
		Repository struct {
			Issue struct {
				Comments struct {
					Nodes    []GithubComment
					PageInfo PageInfo
				} `graphql:"comments(first: 50, after: $cursor)"`
			} `graphql:"issue(number: $number)"`
		} `graphql:"repository(name: $repository, owner: $owner)"`
		RateLimit RateLimit
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repository": githubv4.String(repository),
		"number":     githubv4.Int(number),
		"cursor":     githubv4.String(cursor),
	}

	page := 0
	var comments []GithubComment
	for {
		g.logger.Debug("Querying issue comments", "owner", owner, "repository", repository, "page", page)
		err := g.v4.Query(context.Background(), &query, variables)
		if err != nil {
			return comments, err
		}

		// Process comments
		comments = append(comments, query.Repository.Issue.Comments.Nodes...)

		if !query.Repository.Issue.Comments.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.String(query.Repository.Issue.Comments.PageInfo.EndCursor)
		page++
	}

	return comments, nil
}

func (g *GithubImpl) QueryPullrequests(owner string, repository string, since time.Time, limit int) ([]database.Pullrequest, error) {
	var query struct {
		Repository struct {
			PullRequests struct {
				Nodes      []GithubPullrequest
				PageInfo   PageInfo
				TotalCount int
			} `graphql:"pullRequests(first: 50, after: $cursor, orderBy: { field: UPDATED_AT, direction: DESC })"`
		} `graphql:"repository(name: $repository, owner: $owner)"`
		RateLimit RateLimit
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repository": githubv4.String(repository),
		"cursor":     (*githubv4.String)(nil),
	}

	var ratelimit RateLimit

	done := false
	page := 0
	pullrequests := []database.Pullrequest{}

	for {
		g.logger.Debug("Querying pullrequests", "owner", owner, "repository", repository, "page", page)
		err := g.v4.Query(context.Background(), &query, variables)
		if err != nil {
			return pullrequests, err
		}

		if ratelimit.Cost > ratelimit.Remaining {
			return pullrequests, fmt.Errorf("the query would exceed the current rate limit")
		}

		// Process pullrequests
		for _, p := range query.Repository.PullRequests.Nodes {
			if p.UpdatedAt.Before(since) {
				done = true
				break
			}

			pullrequest := database.Pullrequest{
				ID:                p.ID,
				Owner:             owner,
				Repository:        repository,
				Author:            p.Author.Login,
				AuthorAssociation: p.AuthorAssociation,
				Title:             p.Title,
				Body:              p.Body,
				CreatedAt:         p.CreatedAt,
				PublishedAt:       p.PublishedAt,
				UpdatedAt:         p.UpdatedAt,
				LastEditedAt:      p.LastEditedAt,
				ClosedAt:          p.ClosedAt,
				Closed:            p.Closed,
				MergedAt:          p.MergedAt,
				Merged:            p.Merged,
				Mergeable:         p.Mergeable,
				MergedBy:          p.MergedBy.Login,
				Additions:         p.Additions,
				Deletions:         p.Deletions,
				ChangedFiles:      p.ChangedFiles,
				BaseRefName:       p.BaseRefName,
				HeadRefName:       p.HeadRefName,
				ReviewDecision:    p.ReviewDecision,
				State:             p.State,
				Locked:            p.Locked,
				Labels:            database.StringArray{},
				Reactions:         []database.PullrequestReaction{},
				Comments:          []database.PullrequestComment{},
				Files:             []database.PullrequestFile{},
				Reviews:           []database.PullrequestReview{},
			}

			timelineItems := p.TimelineItems.Nodes
			if len(timelineItems) > 0 {
				// Get the last actor that closed the pullrequest
				pullrequest.ClosedBy = timelineItems[len(timelineItems)-1].PullrequestTimelineItemsConnection.Actor.Login
			}

			// Reactions
			for _, r := range p.ReactionGroups {
				reaction := database.PullrequestReaction{
					Reaction: r.Content,
					Count:    r.Users.TotalCount,
				}
				pullrequest.Reactions = append(pullrequest.Reactions, reaction)
			}

			// Labels
			for _, l := range p.Labels.Nodes {
				pullrequest.Labels = append(pullrequest.Labels, l.Name)
			}

			// Query additional comments
			if p.Comments.PageInfo.HasNextPage {
				g.logger.Debug("Need to query additional comments")
				comments, err := g.QueryIssueComments(owner, repository, p.Number, p.Comments.PageInfo.EndCursor)
				if err != nil {
					return pullrequests, err
				}
				p.Comments.Nodes = append(p.Comments.Nodes, comments...)
			}
			// Comments
			for _, c := range p.Comments.Nodes {
				comment := database.PullrequestComment{
					ID:                c.ID,
					Author:            c.Author.Login,
					AuthorAssociation: c.AuthorAssociation,
					Body:              c.Body,
					CreatedAt:         c.CreatedAt,
					LastEditedAt:      c.LastEditedAt,
					PublishedAt:       c.PublishedAt,
					UpdatedAt:         c.UpdatedAt,
				}

				// Comment reactions
				for _, cr := range c.ReactionGroups {
					reaction := database.PullrequestCommentReaction{
						Content: cr.Content,
						Count:   cr.Users.TotalCount,
					}
					comment.Reactions = append(comment.Reactions, reaction)
				}
				pullrequest.Comments = append(pullrequest.Comments, comment)

				// Reviews
				for _, r := range p.Reviews.Nodes {
					review := database.PullrequestReview{
						Pullrequest:       p.ID,
						Author:            r.Author.Login,
						AuthorAssociation: r.AuthorAssociation,
						Body:              r.Body,
						State:             r.State,
						CreatedAt:         r.CreatedAt,
						PublishedAt:       r.PublishedAt,
						LastEditedAt:      r.LastEditedAt,
						UpdatedAt:         r.UpdatedAt,
						SubmittedAt:       r.SubmittedAt,
					}
					pullrequest.Reviews = append(pullrequest.Reviews, review)
				}

				// Files
				for _, f := range p.Files.Nodes {
					file := database.PullrequestFile{
						Pullrequest: p.ID,
						Path:        f.Path,
						Additions:   f.Additions,
						Deletions:   f.Deletions,
					}
					pullrequest.Files = append(pullrequest.Files, file)
				}
			}

			pullrequests = append(pullrequests, pullrequest)

			if len(pullrequests) == limit {
				return pullrequests, nil
			}
		}

		if done || !query.Repository.PullRequests.PageInfo.HasNextPage {
			break
		}

		ratelimit = query.RateLimit

		variables["cursor"] = githubv4.String(query.Repository.PullRequests.PageInfo.EndCursor)
		page++
	}

	return pullrequests, nil
}

func (g *GithubImpl) QueryPullrequestComments(owner string, repository string, number int, cursor string) ([]GithubComment, error) {
	var query struct {
		Repository struct {
			PullRequest struct {
				Comments struct {
					Nodes    []GithubComment
					PageInfo PageInfo
				} `graphql:"comments(first: 50, after: $cursor)"`
			} `graphql:"pullRequest(number: $number)"`
		} `graphql:"repository(name: $repository, owner: $owner)"`
		RateLimit RateLimit
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repository": githubv4.String(repository),
		"number":     githubv4.Int(number),
		"cursor":     (*githubv4.String)(nil),
	}

	page := 0
	var comments []GithubComment
	for {
		g.logger.Debug("Querying pullrequests comments", "owner", owner, "repository", repository, "page", page)
		err := g.v4.Query(context.Background(), &query, variables)
		if err != nil {
			return comments, err
		}

		// Process comments
		comments = append(comments, query.Repository.PullRequest.Comments.Nodes...)

		if !query.Repository.PullRequest.Comments.PageInfo.HasNextPage {
			break
		}

		variables["cursor"] = githubv4.String(query.Repository.PullRequest.Comments.PageInfo.EndCursor)
		page++
	}

	return comments, nil
}

func (g *GithubImpl) QueryReleases(owner string, repository string, limit int) ([]database.Release, error) {
	var query struct {
		Repository struct {
			Releases struct {
				Nodes      []GithubRelease
				PageInfo   PageInfo
				TotalCount int
			} `graphql:"releases(first: 50, after: $cursor)"`
		} `graphql:"repository(name: $repository, owner: $owner)"`
		RateLimit RateLimit
	}

	variables := map[string]interface{}{
		"owner":      githubv4.String(owner),
		"repository": githubv4.String(repository),
		"cursor":     (*githubv4.String)(nil),
	}

	var ratelimit RateLimit

	done := false
	page := 0
	releases := []database.Release{}

	for {
		g.logger.Debug("Querying releases", "owner", owner, "repository", repository, "page", page)
		err := g.v4.Query(context.Background(), &query, variables)
		if err != nil {
			return releases, err
		}

		if ratelimit.Cost > ratelimit.Remaining {
			return releases, fmt.Errorf("the query would exceed the current rate limit")
		}

		// Process releases
		for _, r := range query.Repository.Releases.Nodes {
			release := database.Release{
				ID:           r.ID,
				Owner:        owner,
				Repository:   repository,
				Name:         r.Name,
				Description:  r.Description,
				CreatedAt:    r.CreatedAt,
				URL:          r.URL,
				IsPrerelease: r.IsPrerelease,
				Tag:          r.TagName,
				Assets:       []database.ReleaseAsset{},
			}

			for _, a := range r.ReleaseAssets.Nodes {
				asset := database.ReleaseAsset{
					ID:         a.ID,
					Release:    r.ID,
					Owner:      owner,
					Repository: repository,
					Name:       a.Name,
					Downloads:  a.DownloadCount,
					Size:       a.Size,
				}

				release.Assets = append(release.Assets, asset)
			}

			releases = append(releases, release)
		}

		if done || !query.Repository.Releases.PageInfo.HasNextPage {
			break
		}

		ratelimit = query.RateLimit

		variables["cursor"] = githubv4.String(query.Repository.Releases.PageInfo.EndCursor)
		page++
	}

	return releases, nil
}

func (g *GithubImpl) QueryMetrics(owner string, repository string) (database.Metrics, error) {
	ctx := context.Background()

	var metrics database.Metrics
	metrics.Owner = owner
	metrics.Repository = repository

	g.logger.Debug("Querying traffic clones", "owner", owner, "repository", repository)
	tc, _, err := g.v3.Repositories.ListTrafficClones(ctx, owner, repository, &githubv3.TrafficBreakdownOptions{})
	if err != nil {
		return metrics, fmt.Errorf("could not query traffic clones: %v+", err)
	}

	metrics.Clones = database.TrafficClones{
		Owner:      owner,
		Repository: repository,
		Count:      tc.GetCount(),
		Uniques:    tc.GetUniques(),
	}

	g.logger.Debug("Querying forks", "owner", owner, "repository", repository)
	forks, _, err := g.v3.Repositories.ListForks(ctx, owner, repository, &githubv3.RepositoryListForksOptions{})
	if err != nil {
		return metrics, fmt.Errorf("could not query forks: %v+", err)
	}

	for _, f := range forks {
		metrics.Forks = append(metrics.Forks, f.GetFullName())
	}

	g.logger.Debug("Querying stargazers", "owner", owner, "repository", repository)
	stars, _, err := g.v3.Activity.ListStargazers(ctx, owner, repository, &githubv3.ListOptions{})
	if err != nil {
		return metrics, fmt.Errorf("could not query stargazers: %v+", err)
	}

	for _, s := range stars {
		metrics.Stars = append(metrics.Stars, s.GetUser().GetLogin())
	}

	g.logger.Debug("Querying watchers", "owner", owner, "repository", repository)
	watches, _, err := g.v3.Activity.ListWatchers(ctx, owner, repository, &githubv3.ListOptions{})
	if err != nil {
		return metrics, fmt.Errorf("could not query watchers: %v+", err)
	}

	for _, w := range watches {
		metrics.Watches = append(metrics.Watches, w.GetLogin())
	}

	g.logger.Debug("Querying traffic views", "owner", owner, "repository", repository)
	tv, _, err := g.v3.Repositories.ListTrafficViews(ctx, owner, repository, &githubv3.TrafficBreakdownOptions{})
	if err != nil {
		return metrics, fmt.Errorf("could not query traffic views: %v+", err)
	}

	metrics.Views = database.TrafficViews{
		Owner:      owner,
		Repository: repository,
		Count:      tv.GetCount(),
		Uniques:    tv.GetUniques(),
	}

	g.logger.Debug("Querying traffic paths", "owner", owner, "repository", repository)
	paths, _, err := g.v3.Repositories.ListTrafficPaths(ctx, owner, repository)
	if err != nil {
		return metrics, fmt.Errorf("could not query traffic paths: %v+", err)
	}

	for _, p := range paths {
		metrics.Paths = append(metrics.Paths, database.TrafficPath{
			Owner:      owner,
			Repository: repository,
			Path:       p.GetPath(),
			Title:      p.GetTitle(),
			Count:      p.GetCount(),
			Uniques:    p.GetUniques(),
		})
	}

	g.logger.Debug("Querying traffic referrers", "owner", owner, "repository", repository)
	referrers, _, err := g.v3.Repositories.ListTrafficReferrers(ctx, owner, repository)
	if err != nil {
		return metrics, fmt.Errorf("could not query traffic referrers: %v+", err)
	}

	for _, r := range referrers {
		metrics.Referrers = append(metrics.Referrers, database.TrafficReferrer{
			Owner:      owner,
			Repository: repository,
			Referrer:   r.GetReferrer(),
			Count:      r.GetCount(),
			Uniques:    r.GetUniques(),
		})
	}

	return metrics, nil
}
