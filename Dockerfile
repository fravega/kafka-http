FROM alpine:edge as base

RUN apk add --no-cache ca-certificates

ENV GOLANG_VERSION 1.9.1

RUN set -eux; \
	apk add --no-cache --virtual .build-deps \
		bash \
		gcc \
		musl-dev \
		openssl \
		go \
	; \
	export \
# set GOROOT_BOOTSTRAP such that we can actually build Go
		GOROOT_BOOTSTRAP="$(go env GOROOT)" \
# ... and set "cross-building" related vars to the installed system's values so that we create a build targeting the proper arch
# (for example, if our build host is GOARCH=amd64, but our build env/image is GOARCH=386, our build needs GOARCH=386)
		GOOS="$(go env GOOS)" \
		GOARCH="$(go env GOARCH)" \
		GO386="$(go env GO386)" \
		GOARM="$(go env GOARM)" \
		GOHOSTOS="$(go env GOHOSTOS)" \
		GOHOSTARCH="$(go env GOHOSTARCH)" \
	; \
	\
	wget -O go.tgz "https://golang.org/dl/go$GOLANG_VERSION.src.tar.gz"; \
	echo 'a84afc9dc7d64fe0fa84d4d735e2ece23831a22117b50dafc75c1484f1cb550e *go.tgz' | sha256sum -c -; \
	tar -C /usr/local -xzf go.tgz; \
	rm go.tgz; \
	\
	cd /usr/local/go/src; \

	./make.bash; \
	\

	apk del .build-deps; \
	\
	export PATH="/usr/local/go/bin:$PATH"; \
	go version

ENV GOPATH /go
ENV PATH $GOPATH/bin:/usr/local/go/bin:$PATH

RUN mkdir -p "$GOPATH/src" "$GOPATH/bin" && chmod -R 777 "$GOPATH"
WORKDIR $GOPATH

#COPY go-wrapper /usr/local/bin/


# See https://hub.docker.com/_/golang/
# See http://blog.wrouesnel.com/articles/Totally%20static%20Go%20builds/
# build stage
FROM base AS build-env

RUN apk update && apk upgrade && \
    apk add --no-cache bash git openssh && \
    apk add librdkafka-dev && apk add pkgconfig && \
    apk add --no-cache  python && \
    apk add --no-cache  gcc g++ && apk add --no-cache  pkgconfig && apk add --no-cache build-base
# RUN git clone https://github.com/edenhill/librdkafka.git

#WORKDIR librdkafka
#RUN ./configure --prefix /usr && make && make install

ENV SRC /go/src/github/fravega/akka-http

ADD main $SRC
COPY ./app /refs
RUN cd $SRC && go get ./... && GOOS=linux go build -o /refs/pricerefs

# final stage
FROM base
RUN apk update && apk upgrade &&  apk add librdkafka-dev
WORKDIR /app
# copy binary into image
COPY --from=build-env /refs /app/
CMD ./pricerefs $BROKER $GROUP $TOPIC
