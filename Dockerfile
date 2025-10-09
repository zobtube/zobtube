FROM alpine:3.22.2

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
