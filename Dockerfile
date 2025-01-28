# BUILDER
FROM cgr.dev/chainguard/go AS builder

COPY ./cmd      /app/cmd
COPY ./internal /app/internal
COPY ./pkg      /app/pkg
COPY ./go.mod   /app/go.mod
COPY ./go.sum   /app/go.sum

RUN cd /app/cmd/bridge; \
    CGO_ENABLED=0 /usr/bin/go build -o /app/bridge .

# FINAL
FROM cgr.dev/chainguard/static

COPY --from=builder /app/bridge /app/bridge

WORKDIR /app

CMD ["/app/bridge"]
