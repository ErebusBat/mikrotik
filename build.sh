#!/usr/bin/env zsh
go build -ldflags "-X main.BUILD_TAG $(git describe --always --dirty)" cmd/bandwidth.go
