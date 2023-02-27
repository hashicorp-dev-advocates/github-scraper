package github

import "time"

type GithubComment struct {
	ID                string
	Author            GithubAuthor
	Body              string
	CreatedAt         time.Time
	LastEditedAt      time.Time
	PublishedAt       time.Time
	UpdatedAt         time.Time
	AuthorAssociation string
	ReactionGroups    []struct {
		Content string
		Users   struct {
			TotalCount int
		}
	}
}

type GithubReaction struct {
	Content string
	Users   struct {
		TotalCount int
	}
}

type GithubLabel struct {
	Name  string
	Color string
}

type GithubAuthor struct {
	Login string
	User  struct {
		Name  string
		Email string
	} `graphql:"... on User"`
}

type GithubActor struct {
	Login string
}

type GithubRelease struct {
	ID            string
	Name          string
	Description   string
	URL           string
	CreatedAt     time.Time
	IsPrerelease  bool
	TagName       string
	ReleaseAssets struct {
		Nodes []GithubReleaseAsset
	} `graphql:"releaseAssets(first: 100)"`
}

type GithubReleaseAsset struct {
	ID            string
	Name          string
	DownloadCount int
	Size          int
}

type GithubIssue struct {
	ID                string
	Number            int
	Author            GithubAuthor
	AuthorAssociation string
	Title             string
	Body              string
	CreatedAt         time.Time
	PublishedAt       time.Time
	UpdatedAt         time.Time
	LastEditedAt      time.Time
	ClosedAt          time.Time
	State             string
	Locked            bool
	Closed            bool
	Comments          struct {
		Nodes    []GithubComment
		PageInfo PageInfo
	} `graphql:"comments(first: 100, after: $cursor)"`
	ReactionGroups []GithubReaction
	Labels         struct {
		Nodes []GithubLabel
	} `graphql:"labels(first: 100)"`
	TimelineItems struct {
		Nodes []struct {
			IssueTimelineItemsConnection struct {
				Actor GithubActor
			} `graphql:"... on ClosedEvent"`
		}
	} `graphql:"timelineItems(itemTypes: CLOSED_EVENT, last: 1)"`
}

type GithubReview struct {
	Author            GithubAuthor
	AuthorAssociation string
	Body              string
	State             string
	CreatedAt         time.Time
	PublishedAt       time.Time
	LastEditedAt      time.Time
	UpdatedAt         time.Time
	SubmittedAt       time.Time
}

type GithubFile struct {
	Additions int
	Deletions int
	Path      string
}

type GithubPullrequest struct {
	ID                string
	Number            int
	Author            GithubAuthor
	AuthorAssociation string
	Title             string
	Body              string
	CreatedAt         time.Time
	ClosedAt          time.Time
	LastEditedAt      time.Time
	MergedAt          time.Time
	UpdatedAt         time.Time
	PublishedAt       time.Time
	Closed            bool
	Merged            bool
	Mergeable         string
	Locked            bool
	Additions         int
	Deletions         int
	ChangedFiles      int
	BaseRefName       string
	HeadRefName       string
	State             string
	ReviewDecision    string
	MergedBy          GithubAuthor
	Comments          struct {
		Nodes    []GithubComment
		PageInfo PageInfo
	} `graphql:"comments(first: 100, after: $cursor)"`
	ReactionGroups []GithubReaction
	Labels         struct {
		Nodes []GithubLabel
	} `graphql:"labels(first: 100)"`
	Reviews struct {
		Nodes []GithubReview
	} `graphql:"reviews(first: 100)"`
	Files struct {
		Nodes []GithubFile
	} `graphql:"files(first: 100)"`
	TimelineItems struct {
		Nodes []struct {
			PullrequestTimelineItemsConnection struct {
				Actor GithubActor
			} `graphql:"... on ClosedEvent"`
		}
	} `graphql:"timelineItems(itemTypes: CLOSED_EVENT, last: 1)"`
}

type PageInfo struct {
	EndCursor   string
	HasNextPage bool
}

type RateLimit struct {
	Cost      int
	Remaining int
	ResetAt   string
}
