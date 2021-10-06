FROM golang:1.13 AS build

ADD . /opt/app
WORKDIR /opt/app
RUN go build ./main.go

EXPOSE 8080

CMD ./main
