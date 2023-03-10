# GitHub data scraper

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
