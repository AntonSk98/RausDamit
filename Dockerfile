FROM golang:1.24.5 AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o raus-damit .

FROM alpine:latest
RUN apk add --no-cache ca-certificates tzdata
WORKDIR /app
COPY --from=builder /app/raus-damit .
COPY config/config.yml /app/config.yml

ENV CONFIG_PATH="/app/config.yml"

ENTRYPOINT ["/app/raus-damit"]