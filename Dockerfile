FROM alpine:3.23.2

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
