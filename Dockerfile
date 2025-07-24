FROM alpine:3.22.1

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
