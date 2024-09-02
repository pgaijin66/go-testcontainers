# Stage 1
FROM golang:1.19.3-alpine AS build_base

ENV CGO_ENABLED=1
ENV GO111MODULE=on
ENV GOOS=linux
ENV GOARCH=amd64

# Set the Current Working Directory inside the container
WORKDIR /src

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

# This is intermediate container hence, we can test for security in here as well
# this will not be in the final package
RUN go install -v golang.org/x/vuln/cmd/govulncheck@latest

RUN CGO_ENABLED=$CGO_ENABLED GOOS=$GOOS GOARCH=$GOARCH go build -ldflags='-s' -o ./out/api ./cmd/api/main.go

# Stage 2
FROM alpine:3.15 as runner

RUN apk update \
  && apk -U upgrade \
  && apk add --no-cache ca-certificates bash gcc \
  && update-ca-certificates --fresh \
  && rm -rf /var/cache/apk/*

RUN addgroup gogroup && adduser -S gouser -u 1000 -G app_user

WORKDIR /app

COPY --chown=gouser:gogroup --from=build_base /src/out/api /app/api

RUN chmod +x api

# This container exposes port 9090 to the outside world
EXPOSE 9090

USER gouser

# Run the binary program produced by `go install`
ENTRYPOINT ["./api"]