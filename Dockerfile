# build env
FROM golang:1.22 AS build-env
COPY go.mod go.sum /src/
WORKDIR /src
RUN go mod download
COPY . .
ARG TARGETOS
ARG TARGETARCH
ARG release=
RUN <<EOR
  VERSION=$(git rev-parse --short HEAD)
  BUILDTIME=$(date -u '+%Y-%m-%dT%H:%M:%SZ')
  RELEASE=$release
  CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -o /app/tx-spammer -ldflags="-s -w -X 'github.com/theQRL/zond-tx-spammer/utils.BuildVersion=${VERSION}' -X 'github.com/theQRL/zond-tx-spammer/utils.BuildRelease=${RELEASE}' -X 'github.com/theQRL/zond-tx-spammer/utils.Buildtime=${BUILDTIME}'" ./cmd/tx-spammer
EOR

# final stage
FROM debian:stable-slim
WORKDIR /app
RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates
RUN update-ca-certificates
ENV PATH="$PATH:/app"
COPY --from=build-env /app/* /app
ENTRYPOINT ["./tx-spammer"]
