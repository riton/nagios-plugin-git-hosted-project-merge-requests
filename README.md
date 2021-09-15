# Nagios plugin to check git hosted project Merge Request

[![Go Report Card](https://goreportcard.com/badge/github.com/riton/nagios-plugin-git-hosted-project-merge-requests)](https://goreportcard.com/report/github.com/riton/nagios-plugin-git-hosted-project-merge-requests)

## Description

This nagios plugin helps project maintainers hosted on Gitlab, Github, etc... to keep track of their _staled_ Merge Request / Pull Requests.

**Important**: Currently only `gitlab` provider is supported

## Build

```
$ CGO_ENABLED=0 go build -ldflags "-s -w -extldflags '-static'" -o check_git_project_merge_requests .
```

## Usage

```
$ check_git_project_merge_requests -h
Checks that a github / gitlab / gitea project has opened merge requests

Usage:
  nagios-plugin-git-hosted-project-merge-requests [flags]

Flags:
      --api-token string                API Token used for authentication
  -c, --config string                   config file (default is /etc/nagios-plugin-git-hosted-project-merge-requests/config.yaml)
      --critical-last-update duration   critical if last-update was that delay ago (default 24h0m0s)
  -d, --debug                           Enable debug
  -p, --git-provider string             git provider can be one of gitlab
  -h, --help                            help for nagios-plugin-git-hosted-project-merge-requests
  -H, --host string                     host to check (API endpoint)
  -P, --project string                  project to check for opened MergeRequests
      --target-branch string            Only consider merge requests with this target-branch (default "master")
  -t, --timeout duration                Global timeout (default 30s)
      --warning-last-update duration    warning if last-update was that delay ago (default 6h0m0s)
```

## Example

### No opened Merge Requests output

```
$ API_TOKEN=XXXXXXX check_git_project_merge_requests -H https://gitlab.com -P "riton/blog" -p gitlab
OK: No opened merge requests | 'total_duration'=0.551541323s;;;; 'opened_merge_requests'=0;;;;
```

### Opened Merge Requests within acceptable range

```
$ API_TOKEN=XXXXXXX check_git_project_merge_requests -H https://gitlab.com -P "riton/blog" -p gitlab
OK: No merge requests too old | 'total_duration'=0.566831438s;;;; 'opened_merge_requests'=1;;;; 'oldest_merge_request'=1319.409961917s;;;;
```

### Opened Merge Requests with too much wait time since last activity

```
$ API_TOKEN=XXXXXXX check_git_project_merge_requests -H https://gitlab.com -P "riton/blog" -p gitlab --warning-last-update 5m --critical-last-update 8m
CRITICAL: Merge request 42 last activity was 21m43.664245342s ago | 'total_duration'=0.795784589s;;;; 'opened_merge_requests'=1;;;; 'oldest_merge_request'=1303.664245342s;;;;
```

## Passing parameters

This project is using [viper](https://github.com/spf13/viper) so any configuration flag can be passed using _environment variables_ or using a configuration file.

A sample configuration file can be found in the [example/config.yaml](example/config.yaml) directory.

## TODO

- [ ] Add support for Github provider
- [ ] Add support for [nagios range](https://nagios-plugins.org/doc/guidelines.html#THRESHOLDFORMAT) definition
