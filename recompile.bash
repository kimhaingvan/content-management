#!/bin/bash
CompileDaemon -log-prefix=false -build="go build ./cmd/content-server/main.go" -command="./main"
