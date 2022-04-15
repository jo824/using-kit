FROM golang:1.17-alpine
RUN apk add --update ca-certificates # Certificates for SSL
RUN echo $(ls -1 /go/bin/)
RUN echo $(ls -1 )
COPY  kit-server-d  .
EXPOSE 8008
ENTRYPOINT kit-server-d

