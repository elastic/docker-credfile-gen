# This Dockerfile is used to cross compile docker-credfile-gen
# to the host machine's OS/ARCH combination inside a Docker container.

FROM golang:1.19 as cache
WORKDIR /cache
COPY go.* /cache/
RUN go mod download

FROM cache
ARG BINARY=docker-credfile-gen
ARG goos=linux
ARG goarch=amd64

WORKDIR /${BINARY}
COPY . /${BINARY}

RUN CGO_ENABLED=0 GOOS=${goos} GOARCH=${goarch} go build -a -o bin/${BINARY} .
