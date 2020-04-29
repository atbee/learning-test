# build
FROM golang:1.14-alpine as builder

RUN apk update && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" >  /etc/timezone \
    && apk del tzdata

WORKDIR /app

COPY . .

RUN go build \
    -ldflags "-X main.buildcommit=$(git rev-parse HEAD) -X main.buildtime=$(date +%Y%m%d.%H%M%S)" \
    -o goapp main.go

# ---------------------------------------------------------

# run
FROM alpine:latest

RUN apk update && apk add tzdata \
    && cp /usr/share/zoneinfo/Asia/Bangkok /etc/localtime \
    && echo "Asia/Bangkok" >  /etc/timezone \
    && apk del tzdata

WORKDIR /app

COPY ./configs ./configs
COPY --from=builder /app/goapp .

CMD ["./goapp"]