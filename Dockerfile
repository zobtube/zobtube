FROM alpine:3.24.2

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
