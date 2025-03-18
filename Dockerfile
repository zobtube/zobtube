FROM scratch

ENTRYPOINT ["/zobtube"]

RUN mkdir /tmp

COPY zobtube /
