FROM golang:1.20 as builder
WORKDIR /app

ARG VERSION
COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN go mod download
COPY ./ ./
RUN make build

FROM debian:bookworm-slim
WORKDIR /app

COPY --from=builder /app/.bin/incident-commander /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/tenant-controller"]
