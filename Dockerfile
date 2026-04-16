FROM alpine:3.23.4

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
