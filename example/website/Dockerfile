# Builder
FROM golang:1.18 AS builder
WORKDIR /go/src/website

COPY . .
RUN go mod init website
RUN go mod tidy

ARG GOOS=linux
ARG GOARCH=amd64
RUN CGO_ENABLED=0 GOOS=$GOOS GOARCH=$GOARCH go build -o /bin/website

# Image
FROM alpine:3.14
COPY --from=builder /bin/website /bin/website
ENTRYPOINT /bin/website
EXPOSE 3000
