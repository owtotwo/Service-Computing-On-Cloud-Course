FROM golang:alpine AS builder
WORKDIR /myapp/
RUN apk --no-cache add git \
  && go get -d -v github.com/unrolled/render \
  && go get -d -v github.com/codegangsta/negroni \
  && go get -d -v github.com/gorilla/mux \
  && go get -d -v github.com/go-sql-driver/mysql \
  && go get -d -v github.com/satori/go.uuid \
  && go get -d -v github.com/spf13/pflag \
  && mkdir -p $GOPATH/src/github.com/owtotwo/ \
  && git clone https://github.com/owtotwo/Service-Computing-On-Cloud-Course.git --branch homework7 \
  --single-branch $GOPATH/src/github.com/owtotwo/Service-Computing-On-Cloud-Course \
  && CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o todolist \
  github.com/owtotwo/Service-Computing-On-Cloud-Course

FROM alpine:latest  
RUN apk --no-cache add ca-certificates \
  && apk --no-cache add mysql-client
WORKDIR /root/
COPY --from=builder /myapp/todolist .
EXPOSE 8080
RUN echo 'while ! mysqladmin ping -h"127.0.0.1:3307" --silent; do sleep 1; done' > waitForMySQL.sh
CMD ["sh", "-c", "sh waitForMySQL.sh && ./todolist"]
