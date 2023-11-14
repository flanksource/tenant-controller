FROM golang:1.20 as builder
WORKDIR /app

COPY go.mod /app/go.mod
COPY go.sum /app/go.sum
RUN go mod download
COPY ./ ./
RUN make build

FROM debian:bookworm-slim
WORKDIR /app

ENV SOPS_VERSION=v3.8.1
ADD --chmod=755 https://github.com/getsops/sops/releases/download/${SOPS_VERSION}/sops-${SOPS_VERSION}.linux.amd64 /usr/local/bin/sops

COPY --from=builder /app/.bin/tenant-controller /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENTRYPOINT ["/app/tenant-controller"]
