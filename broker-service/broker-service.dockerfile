# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY brokerApp /app

RUN apk add curl

CMD ["/app/brokerApp"]


