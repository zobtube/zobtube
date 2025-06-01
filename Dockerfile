FROM alpine:3.22.0

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
