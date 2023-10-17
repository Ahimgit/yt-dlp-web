#!/bin/sh

staticcheck ./cmd/yt-dlp-web
mkdir -p build
env GOOS=windows go build -o ./build/ ./cmd/yt-dlp-web
env GOOS=linux go build -o ./build/ ./cmd/yt-dlp-web