FROM golang:1.9.2-alpine3.6 AS build

EXPOSE 8080
EXPOSE 443
ENV GOPATH /go


WORKDIR /go/src/go-team-room/

RUN mkdir -p /go/src/go-team-room/conf
COPY conf/conf.json /go/src/go-team-room/conf
COPY go-team-room /go/src/go-team-room/
COPY client/dist /go/src/go-team-room/client/dist
COPY client/index.html /go/src/go-team-room/client
CMD ./go-team-room
