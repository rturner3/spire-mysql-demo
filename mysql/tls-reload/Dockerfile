FROM golang:1.21.3 as build-stage

WORKDIR /app

COPY go.mod go.sum ./
COPY *.go ./

RUN go mod download

RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tls-reloader

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /tls-reloader /tls-reloader

ENTRYPOINT ["/tls-reloader"]