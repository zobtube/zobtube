FROM alpine:3.24.4

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
