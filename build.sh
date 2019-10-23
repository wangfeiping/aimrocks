#!/bin/sh

buildDate=`date +"%F %T %z"`
goVersion=`go version`
goVersion=${goVersion#"go version "}

go build -mod=readonly --ldflags "-X main.Version=v0.0.0 \
    -X main.GitCommit=$(git rev-parse HEAD) \
    -X 'main.BuidDate=$buildDate' \
    -X 'main.GoVersion=$goVersion'" \
    -o ./build/aimrockscli ./cmd/aimrockscli/
