#!/bin/zsh
set -e

GOOS=linux GO111MODULE=on go build -o rproxy.linux ./cmd/rproxy
docker build -t rproxy:latest .