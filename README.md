# Heartilly

[![CI](https://github.com/shiimaxx/heartilly/actions/workflows/ci.yaml/badge.svg)](https://github.com/shiimaxx/heartilly/actions/workflows/ci.yaml)
[![Go Report Card](https://goreportcard.com/badge/github.com/shiimaxx/heartilly)](https://goreportcard.com/report/github.com/shiimaxx/heartilly)

A simple uptime monitoring server.

## Features

- Checks to multiple HTTP endpoints continually
- Notifications to Slack channel when detecting error(timeout, error response, ...)

## Usage

```
heartilly -c config.toml
```

## Configuration

```toml
[notification.slack]
token = "token"
channel = "#general"

[[target]]
url = "https://example.com/check"

[[target]]
url = "https://example.com/check_post"
method = "POST"

[[target]]
url = "https://example.com/check_follow_redirect"
follow = true
```

