# build backend
FROM golang:1.13.6-alpine3.11 as backend_builder
RUN apk add --update --no-cache git alpine-sdk upx
ADD . /go
WORKDIR /go/src
RUN GO111MODULE=on go build -ldflags="-s -w"  -o /go/bin/webhookForwarder
RUN upx /go/bin/webhookForwarder

# Build actual image
FROM alpine:3.11
WORKDIR /opt/
COPY --from=backend_builder /go/bin/webhookForwarder .
EXPOSE 8080
ENTRYPOINT [ "/opt/webhookForwarder" ]