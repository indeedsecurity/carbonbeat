# Build stage
ARG GO_VERSION=1.8
ARG PROJECT_PATH=/go/src/github.com/indeedsecurity/carbonbeat
FROM golang:${GO_VERSION}-alpine AS builder
RUN apk --no-cache add ca-certificates tini git gcc musl-dev
RUN go get -u github.com/golang/dep/cmd/dep
WORKDIR /go/src/github.com/indeedsecurity/carbonbeat
COPY ./ ${PROJECT_PATH}
RUN dep ensure
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o carbonbeat \
    && go test $(go list ./... | grep -v /vendor/)

# Production image
FROM alpine:3.6
RUN apk --no-cache add ca-certificates tini
WORKDIR /
COPY --from=builder /go/src/github.com/indeedsecurity/carbonbeat/carbonbeat .
RUN adduser -D -u 69999 -s /usr/sbin/nologin carbonbeat
#USER carbonbeat
ENTRYPOINT ["tini", "-g", "--"]
CMD ["/carbonbeat", "-v", "-e", "-d", "'*'"]
