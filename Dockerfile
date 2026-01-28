FROM alpine:3.23.3

RUN apk add ffmpeg

ENTRYPOINT ["/zobtube"]

COPY zobtube /
