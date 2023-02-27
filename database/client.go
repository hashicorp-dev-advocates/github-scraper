package database

import (
	"time"

	"github.com/hashicorp/go-hclog"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

type Database interface {
	AddIssue(input Issue) (Issue, error)
	AddIssueReaction(input IssueReaction) (IssueReaction, error)
	AddIssueComment(input IssueComment) (IssueComment, error)
	AddIssueCommentReaction(input IssueCommentReaction) (IssueCommentReaction, error)

	AddPullrequest(input Pullrequest) (Pullrequest, error)
	AddPullrequestReaction(input PullrequestReaction) (PullrequestReaction, error)
	AddPullrequestReview(input PullrequestReview) (PullrequestReview, error)
	AddPullrequestFile(input PullrequestFile) (PullrequestFile, error)
	AddPullrequestComment(input PullrequestComment) (PullrequestComment, error)
	AddPullrequestCommentReaction(input PullrequestCommentReaction) (PullrequestCommentReaction, error)

	AddRelease(input Release) (Release, error)
	AddReleaseAsset(input ReleaseAsset) (ReleaseAsset, error)

	AddMetrics(input Metrics) (Metrics, error)

	AddTrafficClones(input TrafficClones) (TrafficClones, error)
	AddTrafficViews(input TrafficViews) (TrafficViews, error)
	AddTrafficReferrer(input TrafficReferrer) (TrafficReferrer, error)
	AddTrafficPath(input TrafficPath) (TrafficPath, error)

	AddMetadata(metadata Metadata) (Metadata, error)
	GetMetadata(owner string, repository string) (Metadata, error)
}

type DatabaseImpl struct {
	client *sqlx.DB
	logger hclog.Logger
}

func New(connection string, DBmaxopenconns int, DBconnmaxlifetime time.Duration, logger hclog.Logger) (Database, error) {
	logger = logger.Named("Database")

	db, err := sqlx.Connect("pgx", connection)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(DBmaxopenconns)
	db.SetConnMaxLifetime(DBconnmaxlifetime)

	return &DatabaseImpl{
		client: db,
		logger: logger,
	}, nil
}
