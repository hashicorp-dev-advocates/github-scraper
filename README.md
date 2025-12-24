# GitHub Data Scraper

A CLI tool for collecting GitHub repository data including releases, pull requests, issues, and metrics. This tool is used by the HashiCorp DevRel team for automated data collection and analysis.

## Overview

This scraper is the core data collection tool used in the DevRel infrastructure ecosystem. It runs as containerized batch jobs on a Nomad cluster and stores data in PostgreSQL for analysis and reporting.

## Related Repositories

### Infrastructure
- **[devrel-infra-main](https://github.com/hashicorp-dev-advocates/devrel-infra-main)** - Provides the Nomad cluster and PostgreSQL database where this tool runs
- **[devrel-github](https://github.com/hashicorp-dev-advocates/devrel-github)** - Contains Nomad job templates that run this scraper tool

### Data Consumers
- **[devrel-slack](https://github.com/hashicorp-dev-advocates/devrel-slack)** - Avocadobot uses collected data for impact tracking
- **[devrel-api](https://github.com/hashicorp-dev-advocates/devrel-api)** - API service that exposes the collected data

## Usage

## Quickstart

Set environment variables for `github` CLI and postgres.

```shell
export GITHUB_TOKEN=
export POSTGRES_USER=user
export POSTGRES_PASSWORD=password
export POSTGRES_DB=database
```

Build binary.

```shell
go build -v -o dist/github .
```

Start local database.

```shell
make
```

Start local database.

```shell
./dist/github issues -f sql -o postgres://user:password@localhost:5432/database hashicorp consul
```

## Releases

Retrieve releases created in the repository:

```shell
github releases hashicorp terraform
```

## Pullrequests

Retrieve pullrequests created in the repository:

```shell
github pullrequests hashicorp terraform
```

Retrieve pullrequests created in the repository that have been updated since `2023-01-31T00:00:00Z`:

```shell
github pullrequests -s "2023-01-31T00:00:00Z" hashicorp terraform
```

## Issues

Retrieve issues created in the repository:

```shell
github issues hashicorp terraform
```

Retrieve issues created in the repository that have been updated since `2023-01-31T00:00:00Z`:

```shell
github issues -s "2023-01-31T00:00:00Z" hashicorp terraform
```

## Metrics

Retrieve traffic metrics about the repository:

```shell
github metrics hashicorp terraform
```

## Output

Output the data as JSON to stdout:

```shell
github issues -f json hashicorp terraform
```

Output the data as JSON to a file at `/tmp/output.json`:

```shell
github issues -f json -o /tmp/output.json hashicorp terraform
```

Output the data as SQL to a database located at `postgres://user:password@host:5432/database`:

```shell
github issues -f sql -o postgres://user:password@host:5432/database hashicorp terraform
```

## Deployment Context

This tool is typically deployed as part of the DevRel infrastructure:

1. **Container Image**: Built and stored in GitHub Container Registry
2. **Nomad Jobs**: Deployed via job templates from [devrel-github](https://github.com/hashicorp-dev-advocates/devrel-github)
3. **Scheduling**: Runs on cron schedules as Nomad batch jobs
4. **Data Storage**: Outputs directly to PostgreSQL database managed by [devrel-infra-main](https://github.com/hashicorp-dev-advocates/devrel-infra-main)

## Integration

The scraper integrates with the broader DevRel ecosystem:

- **Authentication**: Uses GitHub tokens managed in the infrastructure
- **Database**: Connects to PostgreSQL instance from devrel-infra-main
- **Monitoring**: Job status monitored through Nomad cluster
- **Data Pipeline**: Collected data feeds into reporting and analysis tools