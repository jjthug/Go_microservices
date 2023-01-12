# build a tiny docker image
FROM alpine:latest

RUN mkdir /app

COPY mailerServiceApp /app
COPY templates /templates

CMD ["/app/mailerServiceApp"]