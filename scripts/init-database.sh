#!/bin/bash
set -e

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<-EOSQL
-- CREATE USER $POSTGRES_USER WITH PASSWORD '$POSTGRES_PASSWORD';
-- CREATE DATABASE $POSTGRES_DB;

\c $POSTGRES_DB

CREATE TABLE IF NOT EXISTS github_metadata (
  owner VARCHAR(255) NOT NULL,
  repository VARCHAR(255) NOT NULL,
  issues_updated_at TIMESTAMP,
  pullrequests_updated_at TIMESTAMP,
  PRIMARY KEY (owner, repository)
);

--
-- issues
--
CREATE TABLE IF NOT EXISTS github_issues (
  id VARCHAR(255) PRIMARY KEY,
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  number BIGINT,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  author VARCHAR(255) NOT NULL, -- github_users_login
  author_association VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  published_at TIMESTAMP,
  updated_at TIMESTAMP,
  last_edited_at TIMESTAMP,
  closed_at TIMESTAMP,
  state VARCHAR(255) NOT NULL,
  locked BOOLEAN,
  closed BOOLEAN,
  closed_by VARCHAR(255) NOT NULL,
  labels VARCHAR(255)[] NOT NULL
);

CREATE TABLE IF NOT EXISTS github_issues_reactions (
  issue VARCHAR(255), -- github_issues_id
  reaction VARCHAR(255),
  count BIGINT,
  PRIMARY KEY (issue, reaction)
);

CREATE TABLE IF NOT EXISTS github_issues_comments (
  id VARCHAR(255) PRIMARY KEY,
  issue VARCHAR(255) NOT NULL, -- github_issues_id
  author VARCHAR(255) NOT NULL, -- github_users_login
  author_association VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  created_at TIMESTAMP,
  published_at TIMESTAMP,
  updated_at TIMESTAMP,
  last_edited_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS github_issues_comments_reactions (
  comment VARCHAR(255), -- github_issues_comments_id
  reaction VARCHAR(255),
  issue VARCHAR(255) NOT NULL, -- github_issues_id
  count BIGINT,
  PRIMARY KEY (issue, comment, reaction)
);

--
-- pullrequests
--
CREATE TABLE IF NOT EXISTS github_pullrequests (
  id VARCHAR(255) PRIMARY KEY,
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  number BIGINT,
  title TEXT NOT NULL,
  body TEXT NOT NULL,
  author VARCHAR(255) NOT NULL, -- github_users_login
  author_association VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  published_at TIMESTAMP,
  updated_at TIMESTAMP,
  last_edited_at TIMESTAMP,
  closed_at TIMESTAMP,
  state VARCHAR(255) NOT NULL,
  locked BOOLEAN,
  closed BOOLEAN,
  closed_by VARCHAR(255) NOT NULL,
  labels VARCHAR(255)[] NOT NULL,
  merged_at TIMESTAMP,
  merged BOOLEAN,
  mergeable VARCHAR(255) NOT NULL,
  additions BIGINT,
  deletions BIGINT,
  changed_files BIGINT,
  base_ref_name VARCHAR(255) NOT NULL,
  head_ref_name VARCHAR(255) NOT NULL,
  review_decision VARCHAR(255) NOT NULL,
  merged_by VARCHAR(255) NOT NULL -- github_users_login
);

CREATE TABLE IF NOT EXISTS github_pullrequests_reviews (
  pullrequest VARCHAR(255), -- github_pullrequests_id
  body TEXT NOT NULL,
  author VARCHAR(255) NOT NULL, -- github_users_login
  author_association VARCHAR(255) NOT NULL,
  created_at TIMESTAMP,
  published_at TIMESTAMP,
  updated_at TIMESTAMP,
  last_edited_at TIMESTAMP,
  submitted_at TIMESTAMP,
  state VARCHAR(255) NOT NULL,
  PRIMARY KEY (pullrequest, author)
);

CREATE TABLE IF NOT EXISTS github_pullrequests_files (
  pullrequest VARCHAR(255), -- github_pullrequests_id
  additions BIGINT,
  deletions BIGINT,
  path VARCHAR(255) NOT NULL,
  PRIMARY KEY (pullrequest, path)
);

CREATE TABLE IF NOT EXISTS github_pullrequests_reactions (
  pullrequest VARCHAR(255), -- github_pullrequests_id
  reaction VARCHAR(255),
  count BIGINT,
  PRIMARY KEY (pullrequest, reaction)
);

CREATE TABLE IF NOT EXISTS github_pullrequests_comments (
  id VARCHAR(255) PRIMARY KEY,
  pullrequest VARCHAR(255) NOT NULL, -- github_pullrequests_id
  author VARCHAR(255) NOT NULL, -- github_users_login
  author_association VARCHAR(255) NOT NULL,
  body TEXT NOT NULL,
  created_at TIMESTAMP,
  published_at TIMESTAMP,
  updated_at TIMESTAMP,
  last_edited_at TIMESTAMP
);

CREATE TABLE IF NOT EXISTS github_pullrequests_comments_reactions (
  comment VARCHAR(255), -- github_pullrequests_comments_id
  reaction VARCHAR(255),
  pullrequest VARCHAR(255) NOT NULL, -- github_pullrequests_id
  count BIGINT,
  PRIMARY KEY (pullrequest, comment, reaction)
);

--
-- releases
--
CREATE TABLE IF NOT EXISTS github_releases (
  id VARCHAR(255) PRIMARY KEY,
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  name VARCHAR(255) NOT NULL,
  description TEXT NOT NULL,
  url TEXT NOT NULL,
  created_at TIMESTAMP,
  is_prerelease BOOLEAN,
  tag VARCHAR(255) NOT NULL
);

CREATE TABLE IF NOT EXISTS github_releases_assets (
  id VARCHAR(255) PRIMARY KEY,
  release VARCHAR(255) NOT NULL,
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  name VARCHAR(255) NOT NULL,
  downloads BIGINT,
  size BIGINT
);

--
-- metrics
--
CREATE TABLE IF NOT EXISTS github_metrics (
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  stars VARCHAR(255)[] NOT NULL,
  watches VARCHAR(255)[] NOT NULL,
  forks VARCHAR(255)[] NOT NULL,
  PRIMARY KEY (repository, owner)
);

CREATE TABLE IF NOT EXISTS github_metrics_clones (
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  date DATE,
  count BIGINT,
  uniques BIGINT,
  PRIMARY KEY (repository, owner, date)
);

CREATE TABLE IF NOT EXISTS github_metrics_views (
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  date DATE,
  count BIGINT,
  uniques BIGINT,
  PRIMARY KEY (repository, owner, date)
);

CREATE TABLE IF NOT EXISTS github_metrics_paths (
  path VARCHAR(255) NOT NULL,
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  date DATE,
  title VARCHAR(255) NOT NULL,
  count BIGINT,
  uniques BIGINT,
  PRIMARY KEY (repository, owner, path, date)
);

CREATE TABLE IF NOT EXISTS github_metrics_referrers (
  referrer VARCHAR(255) NOT NULL,
  repository VARCHAR(255) NOT NULL, -- github_metadata_repository
  owner VARCHAR(255) NOT NULL, -- github_metadata_owner
  date DATE,
  count BIGINT,
  uniques BIGINT,
  PRIMARY KEY (repository, owner, referrer, date)
);

GRANT ALL ON ALL TABLES IN SCHEMA "public" TO $POSTGRES_USER;
EOSQL