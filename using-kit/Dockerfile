FROM golang:1.17-alpine
RUN apk add --update ca-certificates # Certificates for SSL
COPY  /go/bin/kit-server-d  .
EXPOSE 8008
ENTRYPOINT kit-server-d

