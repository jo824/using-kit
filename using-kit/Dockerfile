FROM golang:1.17-alpine
RUN apk add --update ca-certificates # Certificates for SSL
WORKDIR /server
COPY ./using-kit/kit-server-d  /server/
EXPOSE 8008
CMD ["./kit-server-d"]

