# Builder
FROM golang:1.18 AS builder
WORKDIR /go/src/sync
COPY go.* .
RUN go mod download
COPY . .
ARG GOOS=linux
ARG GOARCH=amd64
RUN CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o /bin/cookie-sync

# Image
FROM alpine:3.14
COPY --from=builder /bin/cookie-sync /bin/cookie-sync
ENTRYPOINT /bin/cookie-sync
