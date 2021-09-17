#!/usr/bin/env bash 
set -xe

# install packages and dependencies
go get github.com/gorilla/mux
go get github.com/gorilla/handlers

# build command
go build -o bin/application application.go
