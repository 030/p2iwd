# P2IWD

[![CI](https://github.com/030/p2iwd/workflows/Go/badge.svg?event=push)](https://github.com/030/p2iwd/actions?query=workflow%3AGo)
![GitHub go.mod Go version](https://img.shields.io/github/go-mod/go-version/030/p2iwd)
[![Go Report Card](https://goreportcard.com/badge/github.com/030/p2iwd)](https://goreportcard.com/report/github.com/030/p2iwd)
[![DevOps SE Questions](https://img.shields.io/stackexchange/devops/t/p2iwd.svg?logo=stackexchange)](https://devops.stackexchange.com/tags/p2iwd)
![Issues](https://img.shields.io/github/issues-raw/030/p2iwd.svg)
![Pull requests](https://img.shields.io/github/issues-pr-raw/030/p2iwd.svg)
![Total downloads](https://img.shields.io/github/downloads/030/p2iwd/total.svg)
![GitHub forks](https://img.shields.io/github/forks/030/p2iwd?label=fork&style=plastic)
![GitHub watchers](https://img.shields.io/github/watchers/030/p2iwd?style=plastic)
![GitHub stars](https://img.shields.io/github/stars/030/p2iwd?style=plastic)
![License](https://img.shields.io/github/license/030/p2iwd.svg)
![Repository Size](https://img.shields.io/github/repo-size/030/p2iwd.svg)
![Contributors](https://img.shields.io/github/contributors/030/p2iwd.svg)
![Commit activity](https://img.shields.io/github/commit-activity/m/030/p2iwd.svg)
![Last commit](https://img.shields.io/github/last-commit/030/p2iwd.svg)
![Release date](https://img.shields.io/github/release-date/030/p2iwd.svg)
![Latest Production Release Version](https://img.shields.io/github/release/030/p2iwd.svg)
[![codecov](https://codecov.io/gh/030/p2iwd/branch/main/graph/badge.svg)](https://codecov.io/gh/030/p2iwd)
[![codebeat badge](https://codebeat.co/badges/72e50a98-d155-4020-a826-89f1a5977249)](https://codebeat.co/projects/github-com-030-p2iwd-main)

Pull and Push Images Without Docker (P2IWD):

- pull an individual image.

## Quickstart

```bash
curl -L https://github.com/030/p2iwd/releases/download/v0.2.0/p2iwd-ubuntu-20.04 -o /tmp/p2iwd-ubuntu-20.04
```

## Configuration

```bash
mkdir ~/.p2iwd
```

~/.p2iwd/config.yml

```bash
---
debug: false
dir: some-dir
host: some-host
push: false
```

~/.p2iwd/creds.yml

```bash
---
pass: some-pass
user: some-user
```

## Usage

```bash
p2iwd
```

## Stargazers over time

[![Stargazers over time](https://starchart.cc/030/p2iwd.svg)](https://starchart.cc/030/p2iwd)
