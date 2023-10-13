FROM golang:1.21-alpine as builder

WORKDIR /usr/local/src

RUN apk --no-cache add bash musl-dev

COPY "./" "./"

RUN go build -o ./bin/app

FROM alpine

COPY --from=builder /usr/local/src/bin/app ./

CMD "./app"