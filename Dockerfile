FROM golang:1.14 AS builder

WORKDIR /src
COPY go.mod go.mod
COPY go.sum go.sum
COPY config.yaml config.yaml
COPY main.go main.go
COPY pkg/ pkg/
COPY vendor/ vendor/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo

FROM alpine:3.11
RUN apk --no-cache add ca-certificates

COPY --from=builder /src/healthcheck /bin/healthcheck
CMD ["/bin/healthcheck"]
