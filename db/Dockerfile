FROM alpine:latest

RUN apk add --no-cache   \
     python3             \
     postgresql

COPY db/entrypoint.sh /entrypoint.sh

ENTRYPOINT /entrypoint.sh
