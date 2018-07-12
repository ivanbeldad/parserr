FROM golang:1.10

LABEL maintainer="Ivan de la Beldad Fernandez <ivandelabeldad@gmail.com>"

ENV GOPATH=/go

ADD . /go/src/parserr

WORKDIR /go/src/parserr

RUN go get ./... && \
    go build -o main .

CMD ["/go/src/parserr/main"]
