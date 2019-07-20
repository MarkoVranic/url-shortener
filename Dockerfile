
FROM golang:latest

RUN go get -v "github.com/go-redis/redis"
RUN go get -v "github.com/gorilla/mux"
RUN go get -v "github.com/speps/go-hashids"
