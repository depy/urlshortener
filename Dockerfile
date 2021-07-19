FROM golang:1.16.6-alpine3.14

RUN mkdir /app
ADD . /app

WORKDIR /app

RUN go build -o urlshortener .

CMD ["/app/urlshortener"]