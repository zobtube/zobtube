FROM alpine:3.21.3

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
