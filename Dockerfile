FROM alpine:3.24.3

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
