FROM scratch

ENTRYPOINT ["/zobtube"]

COPY docker/tmp /tmp
COPY zobtube /
