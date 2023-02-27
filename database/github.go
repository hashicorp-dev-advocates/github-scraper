package database

import (
	"database/sql"
	"time"
)

type Metadata struct {
	Repository            string    `json:"repository" db:"repository"`
	Owner                 string    `json:"owner" db:"owner"`
	IssuesUpdatedAt       time.Time `json:"issues_updated_at" db:"issues_updated_at"`
	PullrequestsUpdatedAt time.Time `json:"pullrequests_updated_at" db:"pullrequests_updated_at"`
}

type Issue struct {
	ID                string          `json:"id" db:"id"`
	Number            int             `json:"number" db:"number"`
	Owner             string          `json:"owner" db:"owner"`
	Repository        string          `json:"repository" db:"repository"`
	Author            string          `json:"author" db:"author"`
	AuthorAssociation string          `json:"author_association" db:"author_association"`
	Title             string          `json:"title" db:"title"`
	Body              string          `json:"body" db:"body"`
	CreatedAt         time.Time       `json:"created_at" db:"created_at"`
	PublishedAt       time.Time       `json:"published_at" db:"published_at"`
	UpdatedAt         time.Time       `json:"updated_at" db:"updated_at"`
	LastEditedAt      time.Time       `json:"last_edited_at" db:"last_edited_at"`
	State             string          `json:"state" db:"state"`
	Locked            bool            `json:"locked" db:"locked"`
	ClosedAt          time.Time       `json:"closed_at" db:"closed_at"`
	Closed            bool            `json:"closed" db:"closed"`
	ClosedBy          string          `json:"closed_by" db:"closed_by"`
	Comments          []IssueComment  `json:"comments" db:"-"`
	Reactions         []IssueReaction `json:"reactions" db:"-"`
	Labels            StringArray     `json:"labels" db:"labels"`
}

type IssueReaction struct {
	Issue    string `json:"-" db:"issue"`
	Reaction string `json:"reaction" db:"reaction"`
	Count    int    `json:"count" db:"count"`
}

type IssueComment struct {
	ID                string                 `json:"id" db:"id"`
	Issue             string                 `json:"-" db:"issue"`
	Author            string                 `json:"author" db:"author"`
	Body              string                 `json:"body" db:"body"`
	CreatedAt         time.Time              `json:"created_at" db:"created_at"`
	LastEditedAt      time.Time              `json:"last_edited_at" db:"last_edited_at"`
	PublishedAt       time.Time              `json:"published_at" db:"published_at"`
	UpdatedAt         time.Time              `json:"updated_at" db:"updated_at"`
	AuthorAssociation string                 `json:"author_association" db:"author_association"`
	Reactions         []IssueCommentReaction `json:"reactions" db:"-"`
}

type IssueCommentReaction struct {
	Issue   string `json:"-" db:"issue"`
	Comment string `json:"-" db:"comment"`
	Content string `json:"reaction" db:"reaction"`
	Count   int    `json:"count" db:"count"`
}

type Pullrequest struct {
	ID                string                `json:"id" db:"id"`
	Number            int                   `json:"number" db:"number"`
	Owner             string                `json:"owner" db:"owner"`
	Repository        string                `json:"repository" db:"repository"`
	Author            string                `json:"author" db:"author"`
	AuthorAssociation string                `json:"author_association" db:"author_association"`
	Title             string                `json:"title" db:"title"`
	Body              string                `json:"body" db:"body"`
	CreatedAt         time.Time             `json:"created_at" db:"created_at"`
	ClosedAt          time.Time             `json:"closed_at" db:"closed_at"`
	LastEditedAt      time.Time             `json:"last_edited_at" db:"last_edited_at"`
	MergedAt          time.Time             `json:"merged_at" db:"merged_at"`
	UpdatedAt         time.Time             `json:"updated_at" db:"updated_at"`
	PublishedAt       time.Time             `json:"published_at" db:"published_at"`
	Closed            bool                  `json:"closed" db:"closed"`
	Merged            bool                  `json:"merged" db:"merged"`
	Mergeable         string                `json:"mergeable" db:"mergeable"`
	Locked            bool                  `json:"locked" db:"locked"`
	Additions         int                   `json:"additions" db:"additions"`
	Deletions         int                   `json:"deletions" db:"deletions"`
	ChangedFiles      int                   `json:"changed_files" db:"changed_files"`
	BaseRefName       string                `json:"base_ref_name" db:"base_ref_name"`
	HeadRefName       string                `json:"head_ref_name" db:"head_ref_name"`
	State             string                `json:"state" db:"state"`
	ReviewDecision    string                `json:"review_decision" db:"review_decision"`
	MergedBy          string                `json:"merged_by" db:"merged_by"`
	ClosedBy          string                `json:"closed_by" db:"closed_by"`
	Labels            StringArray           `json:"labels" db:"labels"`
	Comments          []PullrequestComment  `json:"comments" db:"-"`
	Reactions         []PullrequestReaction `json:"reactions" db:"-"`
	Reviews           []PullrequestReview   `json:"reviews" db:"-"`
	Files             []PullrequestFile     `json:"files" db:"-"`
}

type PullrequestReaction struct {
	Pullrequest string `json:"-" db:"pullrequest"`
	Reaction    string `json:"reaction" db:"reaction"`
	Count       int    `json:"count" db:"count"`
}

type PullrequestComment struct {
	ID                string                       `json:"id" db:"id"`
	Pullrequest       string                       `json:"-" db:"pullrequest"`
	Author            string                       `json:"author" db:"author"`
	Body              string                       `json:"body" db:"body"`
	CreatedAt         time.Time                    `json:"created_at" db:"created_at"`
	LastEditedAt      time.Time                    `json:"last_edited_at" db:"last_edited_at"`
	PublishedAt       time.Time                    `json:"published_at" db:"published_at"`
	UpdatedAt         time.Time                    `json:"updated_at" db:"updated_at"`
	AuthorAssociation string                       `json:"author_association" db:"author_association"`
	Reactions         []PullrequestCommentReaction `json:"reactions" db:"-"`
}

type PullrequestCommentReaction struct {
	Pullrequest string `json:"-" db:"pullrequest"`
	Comment     string `json:"-" db:"comment"`
	Content     string `json:"reaction" db:"reaction"`
	Count       int    `json:"count" db:"count"`
}

type PullrequestReview struct {
	Pullrequest       string    `json:"-" db:"pullrequest"`
	Author            string    `json:"author" db:"author"`
	AuthorAssociation string    `json:"author_association" db:"author_association"`
	Body              string    `json:"body" db:"body"`
	State             string    `json:"state" db:"state"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	PublishedAt       time.Time `json:"published_at" db:"published_at"`
	LastEditedAt      time.Time `json:"last_edited_at" db:"last_edited_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
	SubmittedAt       time.Time `json:"submitted_at" db:"submitted_at"`
}

type PullrequestFile struct {
	Pullrequest string `json:"-" db:"pullrequest"`
	Additions   int    `json:"additions" db:"additions"`
	Deletions   int    `json:"deletions" db:"deletions"`
	Path        string `json:"path" db:"path"`
}

type Release struct {
	ID           string         `json:"id" db:"id"`
	Owner        string         `json:"owner" db:"owner"`
	Repository   string         `json:"repository" db:"repository"`
	Name         string         `json:"name" db:"name"`
	Description  string         `json:"description" db:"description"`
	URL          string         `json:"url" db:"url"`
	CreatedAt    time.Time      `json:"created_at" db:"created_at"`
	IsPrerelease bool           `json:"is_prerelease" db:"is_prerelease"`
	Tag          string         `json:"tag" db:"tag"`
	Assets       []ReleaseAsset `json:"assets" db:"-"`
}

type ReleaseAsset struct {
	ID         string `json:"id" db:"id"`
	Release    string `json:"-" db:"release"`
	Owner      string `json:"-" db:"owner"`
	Repository string `json:"-" db:"repository"`
	Name       string `json:"name" db:"name"`
	Downloads  int    `json:"downloads" db:"downloads"`
	Size       int    `json:"size" db:"size"`
}

type Metrics struct {
	Owner      string            `json:"owner" db:"owner"`
	Repository string            `json:"repository" db:"repository"`
	Forks      StringArray       `json:"forks" db:"forks"`
	Watches    StringArray       `json:"watches" db:"watches"`
	Stars      StringArray       `json:"stars" db:"stars"`
	Clones     TrafficClones     `json:"clones" db:"-"`
	Views      TrafficViews      `json:"views" db:"-"`
	Paths      []TrafficPath     `json:"paths" db:"-"`
	Referrers  []TrafficReferrer `json:"referrers" db:"-"`
}

type TrafficClones struct {
	Owner      string `json:"-" db:"owner"`
	Repository string `json:"-" db:"repository"`
	Count      int    `json:"count" db:"count"`
	Uniques    int    `json:"uniques" db:"uniques"`
}

type TrafficViews struct {
	Owner      string `json:"-" db:"owner"`
	Repository string `json:"-" db:"repository"`
	Count      int    `json:"count" db:"count"`
	Uniques    int    `json:"uniques" db:"uniques"`
}

type TrafficPath struct {
	Path       string `json:"path" db:"path"`
	Owner      string `json:"-" db:"owner"`
	Repository string `json:"-" db:"repository"`
	Title      string `json:"title" db:"title"`
	Count      int    `json:"count" db:"count"`
	Uniques    int    `json:"uniques" db:"uniques"`
}

type TrafficReferrer struct {
	Referrer   string `json:"referrer" db:"referrer"`
	Owner      string `json:"-" db:"owner"`
	Repository string `json:"-" db:"repository"`
	Count      int    `json:"count" db:"count"`
	Uniques    int    `json:"uniques" db:"uniques"`
}

// Metadata
func (db *DatabaseImpl) AddMetadata(input Metadata) (Metadata, error) {
	var metadata Metadata
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_metadata (
			owner,
			repository,
			issues_updated_at,
			pullrequests_updated_at
		)
		VALUES (
			:owner,
			:repository,
			:issues_updated_at,
			:pullrequests_updated_at
		)
		ON CONFLICT (owner, repository) DO UPDATE 
		SET 
			issues_updated_at = EXCLUDED.issues_updated_at, 
			pullrequests_updated_at = EXCLUDED.pullrequests_updated_at
		RETURNING *`)
	if err != nil {
		return metadata, err
	}

	err = query.Get(&metadata, input)
	if err != nil {
		return metadata, err
	}

	query.Close()

	return metadata, nil
}

func (db *DatabaseImpl) GetMetadata(owner string, repository string) (Metadata, error) {
	var metadata Metadata

	params := map[string]interface{}{
		"owner":      owner,
		"repository": repository,
	}

	query, err := db.client.PrepareNamed(
		`SELECT * FROM github_metadata WHERE owner = :owner AND repository = :repository`)
	if err != nil {
		return metadata, err
	}

	err = query.Get(&metadata, params)
	if err != nil {
		if err == sql.ErrNoRows {
			metadata.Owner = owner
			metadata.Repository = repository
			return metadata, nil
		}

		return metadata, err
	}

	return metadata, nil
}

// Issues
func (db *DatabaseImpl) AddIssue(input Issue) (Issue, error) {
	var issue Issue
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_issues (
			id,
			repository,
			owner,
			number,
			title,
			body,
			author,
			author_association,
			created_at,
			published_at,
			updated_at,
			last_edited_at,
			closed_at,
			state,
			locked,
			closed,
			labels,
			closed_by
		)
		VALUES (
			:id,
			:repository,
			:owner,
			:number,
			:title,
			:body,
			:author,
			:author_association,
			:created_at,
			:published_at,
			:updated_at,
			:last_edited_at,
			:closed_at,
			:state,
			:locked,
			:closed,
			:labels,
			:closed_by
		)
		ON CONFLICT (id) DO UPDATE 
		SET
			title = EXCLUDED.title,
			body = EXCLUDED.body,
			author = EXCLUDED.author,
			author_association = EXCLUDED.author_association,
			created_at = EXCLUDED.created_at,
			published_at = EXCLUDED.published_at,
			updated_at = EXCLUDED.updated_at,
			last_edited_at = EXCLUDED.last_edited_at,
			closed_at = EXCLUDED.closed_at,
			state = EXCLUDED.state,
			locked = EXCLUDED.locked,
			closed = EXCLUDED.closed,
			labels = EXCLUDED.labels,
			closed_by = EXCLUDED.closed_by
		RETURNING *`)
	if err != nil {
		return issue, err
	}

	err = query.Get(&issue, input)
	if err != nil {
		return issue, err
	}

	query.Close()

	return issue, nil
}

func (db *DatabaseImpl) AddIssueReaction(input IssueReaction) (IssueReaction, error) {
	var reaction IssueReaction
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_issues_reactions (
			issue,
			reaction,
			count
		)
		VALUES (
			:issue,
			:reaction,
			:count
		)
		ON CONFLICT (issue, reaction) DO UPDATE SET count = EXCLUDED.count
		RETURNING *`)
	if err != nil {
		return reaction, err
	}

	err = query.Get(&reaction, input)
	if err != nil {
		return reaction, err
	}

	query.Close()

	return reaction, nil
}

func (db *DatabaseImpl) AddIssueComment(input IssueComment) (IssueComment, error) {
	var comment IssueComment
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_issues_comments (
			id,
			issue,
			author,
			author_association,
			body,
			created_at,
			published_at,
			updated_at,
			last_edited_at
		)
		VALUES (
			:id,
			:issue,
			:author,
			:author_association,
			:body,
			:created_at,
			:published_at,
			:updated_at,
			:last_edited_at
		)
		ON CONFLICT (id) DO UPDATE 
		SET 
			issue = EXCLUDED.issue, 
			author = EXCLUDED.author, 
			author_association = EXCLUDED.author_association,
			body = EXCLUDED.body, 
			created_at = EXCLUDED.created_at, 
			published_at = EXCLUDED.published_at, 
			updated_at = EXCLUDED.updated_at, 
			last_edited_at = EXCLUDED.last_edited_at
		RETURNING *`)
	if err != nil {
		return comment, err
	}

	err = query.Get(&comment, input)
	if err != nil {
		return comment, err
	}

	query.Close()

	return comment, nil
}

func (db *DatabaseImpl) AddIssueCommentReaction(input IssueCommentReaction) (IssueCommentReaction, error) {
	var reaction IssueCommentReaction
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_issues_comments_reactions (
			issue,
			comment,
			reaction,
			count
		)
		VALUES (
			:issue,
			:comment,
			:reaction,
			:count
		)
		ON CONFLICT (issue, comment, reaction) DO UPDATE SET count = EXCLUDED.count
		RETURNING *`)
	if err != nil {
		return reaction, err
	}

	err = query.Get(&reaction, input)
	if err != nil {
		return reaction, err
	}

	query.Close()

	return reaction, nil
}

// Pullrequests
func (db *DatabaseImpl) AddPullrequest(input Pullrequest) (Pullrequest, error) {
	var pr Pullrequest
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_pullrequests (
			id,
			repository,
			owner,
			number,
			title,
			body,
			author,
			author_association,
			created_at,
			published_at,
			updated_at,
			last_edited_at,
			closed_at,
			state,
			locked,
			closed,
			labels,
			merged_at,
			merged,
			mergeable,
			additions,
			deletions,
			changed_files,
			base_ref_name,
			head_ref_name,
			review_decision,
			merged_by,
			closed_by
		)
		VALUES (
			:id,
			:repository,
			:owner,
			:number,
			:title,
			:body,
			:author,
			:author_association,
			:created_at,
			:published_at,
			:updated_at,
			:last_edited_at,
			:closed_at,
			:state,
			:locked,
			:closed,
			:labels,
			:merged_at,
			:merged,
			:mergeable,
			:additions,
			:deletions,
			:changed_files,
			:base_ref_name,
			:head_ref_name,
			:review_decision,
			:merged_by,
			:closed_by
		)
		ON CONFLICT (id) DO UPDATE 
		SET
			title = EXCLUDED.title,
			body = EXCLUDED.body,
			author = EXCLUDED.author,
			author_association = EXCLUDED.author_association,
			created_at = EXCLUDED.created_at,
			published_at = EXCLUDED.published_at,
			updated_at = EXCLUDED.updated_at,
			last_edited_at = EXCLUDED.last_edited_at,
			closed_at = EXCLUDED.closed_at,
			state = EXCLUDED.state,
			locked = EXCLUDED.locked,
			closed = EXCLUDED.closed,
			labels = EXCLUDED.labels,
			merged_at = EXCLUDED.merged_at,
			merged = EXCLUDED.merged,
			mergeable = EXCLUDED.mergeable,
			additions = EXCLUDED.additions,
			deletions = EXCLUDED.deletions,
			changed_files = EXCLUDED.changed_files,
			base_ref_name = EXCLUDED.base_ref_name,
			head_ref_name = EXCLUDED.head_ref_name,
			review_decision = EXCLUDED.review_decision,
			merged_by = EXCLUDED.merged_by,
			closed_by = EXCLUDED.closed_by
		RETURNING *`)
	if err != nil {
		return pr, err
	}

	err = query.Get(&pr, input)
	if err != nil {
		return pr, err
	}

	query.Close()

	return pr, nil
}

func (db *DatabaseImpl) AddPullrequestReaction(input PullrequestReaction) (PullrequestReaction, error) {
	var reaction PullrequestReaction
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_pullrequests_reactions (
			pullrequest,
			reaction,
			count
		)
		VALUES (
			:pullrequest,
			:reaction,
			:count
		)
		ON CONFLICT (pullrequest, reaction) DO UPDATE SET count = EXCLUDED.count
		RETURNING *`)
	if err != nil {
		return reaction, err
	}

	err = query.Get(&reaction, input)
	if err != nil {
		return reaction, err
	}

	query.Close()

	return reaction, nil
}

func (db *DatabaseImpl) AddPullrequestReview(input PullrequestReview) (PullrequestReview, error) {
	var review PullrequestReview
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_pullrequests_reviews (
			pullrequest,
			body,
			author,
			author_association,
			created_at,
			published_at,
			updated_at,
			last_edited_at,
			submitted_at,
  		state
		)
		VALUES (
			:pullrequest,
			:body,
			:author,
			:author_association,
			:created_at,
			:published_at,
			:updated_at,
			:last_edited_at,
			:submitted_at,
  		:state
		)
		ON CONFLICT (pullrequest, author) DO UPDATE 
		SET 
			pullrequest = EXCLUDED.pullrequest,
			body = EXCLUDED.body,
			author_association = EXCLUDED.author_association,
			created_at = EXCLUDED.created_at,
			published_at = EXCLUDED.published_at,
			updated_at = EXCLUDED.updated_at,
			last_edited_at = EXCLUDED.last_edited_at,
			submitted_at = EXCLUDED.submitted_at,
			state = EXCLUDED.state
		RETURNING *`)
	if err != nil {
		return review, err
	}

	err = query.Get(&review, input)
	if err != nil {
		return review, err
	}

	query.Close()

	return review, nil
}

func (db *DatabaseImpl) AddPullrequestFile(input PullrequestFile) (PullrequestFile, error) {
	var file PullrequestFile
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_pullrequests_files (
			pullrequest,
			additions,
			deletions,
			path
		)
		VALUES (
			:pullrequest,
			:additions,
			:deletions,
			:path
		)
		ON CONFLICT (pullrequest, path) DO UPDATE 
		SET
			pullrequest = EXCLUDED.pullrequest,
			additions = EXCLUDED.additions,
			deletions = EXCLUDED.deletions,
			path = EXCLUDED.path
		RETURNING *`)
	if err != nil {
		return file, err
	}

	err = query.Get(&file, input)
	if err != nil {
		return file, err
	}

	query.Close()

	return file, nil
}

func (db *DatabaseImpl) AddPullrequestComment(input PullrequestComment) (PullrequestComment, error) {
	var comment PullrequestComment
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_pullrequests_comments (
			id,
			pullrequest,
			author,
			author_association,
			body,
			created_at,
			published_at,
			updated_at,
			last_edited_at
		)
		VALUES (
			:id,
			:pullrequest,
			:author,
			:author_association,
			:body,
			:created_at,
			:published_at,
			:updated_at,
			:last_edited_at
		)
		ON CONFLICT (id) DO UPDATE 
		SET 
			pullrequest = EXCLUDED.pullrequest, 
			author = EXCLUDED.author, 
			author_association = EXCLUDED.author_association,
			body = EXCLUDED.body, 
			created_at = EXCLUDED.created_at, 
			published_at = EXCLUDED.published_at, 
			updated_at = EXCLUDED.updated_at, 
			last_edited_at = EXCLUDED.last_edited_at
		RETURNING *`)
	if err != nil {
		return comment, err
	}

	err = query.Get(&comment, input)
	if err != nil {
		return comment, err
	}

	query.Close()

	return comment, nil
}

func (db *DatabaseImpl) AddPullrequestCommentReaction(input PullrequestCommentReaction) (PullrequestCommentReaction, error) {
	var reaction PullrequestCommentReaction
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_pullrequests_comments_reactions (
			pullrequest,
			comment,
			reaction,
			count
		)
		VALUES (
			:pullrequest,
			:comment,
			:reaction,
			:count
		)
		ON CONFLICT (pullrequest, comment, reaction) DO UPDATE SET count = EXCLUDED.count
		RETURNING *`)
	if err != nil {
		return reaction, err
	}

	err = query.Get(&reaction, input)
	if err != nil {
		return reaction, err
	}

	query.Close()

	return reaction, nil
}

func (db *DatabaseImpl) AddRelease(input Release) (Release, error) {
	var release Release
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_releases (
			id,
			owner,
			repository,
			name,
			description,
			url,
			created_at,
			is_prerelease,
			tag
		)
		VALUES (
			:id,
			:owner,
			:repository,
			:name,
			:description,
			:url,
			:created_at,
			:is_prerelease,
			:tag
		)
		ON CONFLICT (id) DO UPDATE 
		SET 
			name = EXCLUDED.name,
			description = EXCLUDED.description,
			url = EXCLUDED.url,
			created_at = EXCLUDED.created_at,
			is_prerelease = EXCLUDED.is_prerelease,
			tag = EXCLUDED.tag
		RETURNING *`)
	if err != nil {
		return release, err
	}

	err = query.Get(&release, input)
	if err != nil {
		return release, err
	}

	query.Close()

	return release, nil
}

func (db *DatabaseImpl) AddReleaseAsset(input ReleaseAsset) (ReleaseAsset, error) {
	var asset ReleaseAsset
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_releases_assets (
			id,
			release,
			owner,
			repository,
			name,
			downloads, 
			size
		)
		VALUES (
			:id,
			:release,
			:owner,
			:repository,
			:name,
			:downloads, 
			:size
		)
		ON CONFLICT (id) DO UPDATE 
		SET 
			downloads = EXCLUDED.downloads
		RETURNING *`)
	if err != nil {
		return asset, err
	}

	err = query.Get(&asset, input)
	if err != nil {
		return asset, err
	}

	query.Close()

	return asset, nil
}

func (db *DatabaseImpl) AddMetrics(input Metrics) (Metrics, error) {
	var metrics Metrics
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_metrics (
			owner,
			repository,
			forks,
			watches,
			stars
		)
		VALUES (
			:owner,
			:repository,
			:forks,
			:watches,
			:stars
		)
		ON CONFLICT (owner, repository) DO UPDATE
		SET 
			forks = EXCLUDED.forks,
			watches = EXCLUDED.watches,
			stars = EXCLUDED.stars
		RETURNING *`)
	if err != nil {
		return metrics, err
	}

	err = query.Get(&metrics, input)
	if err != nil {
		return metrics, err
	}

	query.Close()

	return metrics, nil
}

func (db *DatabaseImpl) AddTrafficClones(input TrafficClones) (TrafficClones, error) {
	var clones TrafficClones
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_metrics_clones (
			owner,
			repository,
			count,
			uniques
		)
		VALUES (
			:owner,
			:repository,
			:count,
			:uniques
		)
		ON CONFLICT (owner, repository) DO UPDATE
		SET 
			count = EXCLUDED.count,
			uniques = EXCLUDED.uniques
		RETURNING *`)
	if err != nil {
		return clones, err
	}

	err = query.Get(&clones, input)
	if err != nil {
		return clones, err
	}

	query.Close()

	return clones, nil
}

func (db *DatabaseImpl) AddTrafficViews(input TrafficViews) (TrafficViews, error) {
	var views TrafficViews
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_metrics_views (
			owner,
			repository,
			count,
			uniques
		)
		VALUES (
			:owner,
			:repository,
			:count,
			:uniques
		)
		ON CONFLICT (owner, repository) DO UPDATE
		SET 
			count = EXCLUDED.count,
			uniques = EXCLUDED.uniques
		RETURNING *`)
	if err != nil {
		return views, err
	}

	err = query.Get(&views, input)
	if err != nil {
		return views, err
	}

	query.Close()

	return views, nil
}

func (db *DatabaseImpl) AddTrafficReferrer(input TrafficReferrer) (TrafficReferrer, error) {
	var referrer TrafficReferrer
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_metrics_referrers (
			owner,
			repository,
			referrer,
			count,
			uniques
		)
		VALUES (
			:owner,
			:repository,
			:referrer,
			:count,
			:uniques
		)
		ON CONFLICT (owner, repository, referrer) DO UPDATE
		SET 
			count = EXCLUDED.count,
			uniques = EXCLUDED.uniques
		RETURNING *`)
	if err != nil {
		return referrer, err
	}

	err = query.Get(&referrer, input)
	if err != nil {
		return referrer, err
	}

	query.Close()

	return referrer, nil
}

func (db *DatabaseImpl) AddTrafficPath(input TrafficPath) (TrafficPath, error) {
	var path TrafficPath
	query, err := db.client.PrepareNamed(
		`INSERT INTO github_metrics_paths (
			owner,
			repository,
			path,
			title,
			count,
			uniques
		)
		VALUES (
			:owner,
			:repository,
			:path,
			:title,
			:count,
			:uniques
		)
		ON CONFLICT (owner, repository, path) DO UPDATE
		SET 
			title = EXCLUDED.title,
			count = EXCLUDED.count,
			uniques = EXCLUDED.uniques
		RETURNING *`)
	if err != nil {
		return path, err
	}

	err = query.Get(&path, input)
	if err != nil {
		return path, err
	}

	query.Close()

	return path, nil
}
