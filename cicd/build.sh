#!/bin/sh

set -e
set -x

export CGO_ENABLED=0
export GOOS=linux
export GOARCH=amd64
export SSHPASS=$SSH_PASS

cd animated-robot
go get -v ./...
go test -v ./...
cd cmd
go build -v -a -installsuffix cgo -o app .
sshpass -v -e scp -o "StrictHostKeyChecking no" app $SSH_USER@$SSH_SERVER:$SSH_FOLDER
