#!/usr/bin/env bash

docker run --rm -ti -v "$PWD":/go/src/github.com/fravega/kafka-http -w /go/src/github.com/fravega/kafka-http golang:1.8 go get -t -u ./... '&&' go test -i -v -x .
