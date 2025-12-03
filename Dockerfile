FROM alpine:3.23.0

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
