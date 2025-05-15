FROM golang:alpine

WORKDIR /tmp/build

COPY go.mod main.go ./
RUN go build
RUN mv fileenv /usr/bin

WORKDIR /
RUN rm -rf /tmp/build
