# step 1 build the executable binary
FROM golang:alpine As builder

RUN apk update && apk add --no-cache git

WORKDIR $GOPATH/src/go_auth_proxy
COPY src/. .

RUN go get -d -v
RUN go build -o /go/bin/go_micmute_server

# step 2 build a minimal image from scratch (just binary + html/js)
FROM scratch
COPY --from=builder /go/bin/go_micmute_server /go/bin/go_micmute_server
COPY src/public/. /var/www/go_micmute_server/

EXPOSE 3003

ENTRYPOINT [ "/go/bin/go_micmute_server"]