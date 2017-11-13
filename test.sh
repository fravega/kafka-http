#!/usr/bin/env bash

SRC=/go/src/github.com/fravega/kafka-http

docker run --rm -ti \
    -v "$PWD":${SRC} \
    -w ${SRC} \
    golang:1.9.1 \
    /bin/bash -c 'go get -t -v ./... && go test -i -v -x .'
