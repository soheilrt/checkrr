FROM golang:1.22-bullseye AS builder

ENV GO111MODULE on

WORKDIR $GOPATH/src/github.com/hellofresh/wms-adapter-service

COPY blockerr .

RUN go build -o /usr/bin/checkrr ./

FROM ubuntu:jammy

COPY --from=builder /usr/bin/checkrr /usr/bin/checkrr

CMD ["checkrr"]

