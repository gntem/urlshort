# Stage 1 (to create a "build" image, ~850MB)
FROM golang:1.11.2 AS builder
RUN go version

COPY . /go/src/github.com/gntem/urlshort/

WORKDIR /go/src/github.com/gntem/urlshort/

RUN set -x && \
    go get github.com/golang/dep/cmd/dep && \
    dep ensure -v

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o app .

# Stage 2 (to create a downsized "container executable", ~7MB)
FROM scratch
WORKDIR /root/
COPY --from=builder /go/src/github.com/gntem/urlshort/app .

EXPOSE 8000
ENTRYPOINT ["./app"]