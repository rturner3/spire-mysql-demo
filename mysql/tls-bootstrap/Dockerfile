FROM golang:1.19 as build-stage

RUN apt-get update && apt-get install -y ca-certificates openssl

ARG cert_location=/usr/local/share/ca-certificates

# Get certificate from "github.com" and "proxy.golang.org".
# This is required for go mod download to run successfully
RUN openssl s_client -showcerts -connect github.com:443 </dev/null 2>/dev/null|openssl x509 -outform PEM > ${cert_location}/github.crt
RUN openssl s_client -showcerts -connect proxy.golang.org:443 </dev/null 2>/dev/null|openssl x509 -outform PEM >  ${cert_location}/proxy.golang.crt
# Update certificates
RUN update-ca-certificates

WORKDIR /app

COPY go.mod go.sum ./
COPY *.go ./

RUN go mod download

RUN GO111MODULE="on" CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /tls-bootstrap

# Deploy the application binary into a lean image
FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /tls-bootstrap /tls-bootstrap

ENTRYPOINT ["/tls-bootstrap"]