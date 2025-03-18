FROM busybox AS build-env

RUN mkdir /build/tmp

FROM scratch

ENTRYPOINT ["/zobtube"]

COPY --from=build-env /build/tmp /tmp
COPY zobtube /
