FROM golang:1.17-alpine
RUN apk add --update ca-certificates # Certificates for SSL
COPY  /go/bin/kit-server-d  /go/bin/
EXPOSE 8008
ENTRYPOINT /go/bin/kit-server-d

