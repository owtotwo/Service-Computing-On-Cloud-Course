FROM golang:alpine AS builder
WORKDIR /myapp/
RUN apk --no-cache add git \
  && go get -d -v github.com/unrolled/render \
  && go get -d -v github.com/codegangsta/negroni \
  && go get -d -v github.com/gorilla/mux \
  && go get -d -v github.com/go-sql-driver/mysql \
  && go get -d -v github.com/satori/go.uuid \
  && mkdir -p $GOPATH/src/github.com/owtotwo/ \
  && git clone https://github.com/owtotwo/Service-Computing-On-Cloud-Course.git --branch homework7 --single-branch $GOPATH/src/github.com/owtotwo/ \
  && CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o todolist \
  github.com/owtotwo/Service-Computing-On-Cloud-Course

FROM alpine:latest  
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /myapp/todolist .
EXPOSE 8080
CMD ["./app"]
