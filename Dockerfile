FROM alpine:3.24.1

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
