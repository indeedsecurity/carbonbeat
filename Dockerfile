# Build stage
FROM golang:1.15-buster as builder
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=`go env GOHOSTOS` GOARCH=`go env GOHOSTARCH` go build -o carbonbeat \
    && go test $(go list ./... | grep -v /vendor/)

# Production image
FROM alpine:3.12
RUN apk --no-cache add ca-certificates tini
WORKDIR /
COPY --from=builder /build/carbonbeat .
RUN adduser -D -u 69999 -s /usr/sbin/nologin carbonbeat
#USER carbonbeat
ENTRYPOINT ["tini", "-g", "--"]
CMD ["/carbonbeat", "-v", "-e", "-d", "'*'"]
