# syntax=docker/dockerfile:1.7

# ---- Build Stage ----
ARG BUILDPLATFORM=linux/amd64
FROM --platform=$BUILDPLATFORM golang:1.26-alpine AS builder

ARG TARGETOS
ARG TARGETARCH
ARG TARGETVARIANT

# Pure-Go build: no CGO, no gcc needed (modernc.org/sqlite replaces go-sqlite3)
RUN apk add --no-cache git ca-certificates tzdata

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN set -eux; \
    GOOS="${TARGETOS:-linux}"; \
    GOARCH="${TARGETARCH:-amd64}"; \
    if [ "$GOARCH" = "arm" ] && [ -n "${TARGETVARIANT}" ]; then \
        GOARM="${TARGETVARIANT#v}"; \
    fi; \
    CGO_ENABLED=0 GOOS="$GOOS" GOARCH="$GOARCH" ${GOARM:+GOARM=$GOARM} \
    go build -ldflags="-s -w" -o /picobot ./cmd/picobot

# ---- Runtime Stage ----
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /picobot /usr/local/bin/picobot

# Distroless 'nonroot' image runs as uid 65534 by default.
# Home directory for the nonroot user is /home/nonroot.
ENV PICOBOT_HOME=/home/nonroot/.picobot

ENTRYPOINT ["picobot"]
CMD ["gateway"]
