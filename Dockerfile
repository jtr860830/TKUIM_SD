FROM golang:1.11
MAINTAINER Jezrien Hsieh

WORKDIR /app
ADD . /app

RUN go build
EXPOSE 8080
CMD /app/SD-Backend
