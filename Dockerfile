FROM golang:1.14 AS builder

WORKDIR /src
COPY go.mod go.mod
COPY go.sum go.sum
COPY main.go main.go
COPY config.yaml config.yaml
COPY pkg/ pkg/
COPY vendor/ vendor/
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=vendor -a -installsuffix cgo

FROM alpine:3.11
RUN apk --no-cache add ca-certificates

COPY --from=builder /src/healthcheck /bin/healthcheck
COPY --from=builder /src/config.yaml config.yaml
CMD ["/bin/healthcheck"]
