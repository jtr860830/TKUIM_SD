FROM golang:1.11
LABEL maintainer="Jezrien Hsieh"

WORKDIR /app
ADD . /app

RUN go build
EXPOSE 8080
ENV DBMS=mysql
ENV DBLC=root:password@tcp(database:3306)/sd?charset=utf8&parseTime=True&loc=Local
CMD /app/SD-Backend
