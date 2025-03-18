FROM busybox AS build-env

RUN mkdir -p /build/tmp

FROM scratch

ENTRYPOINT ["/zobtube"]

COPY --from=build-env /build/tmp /tmp
COPY zobtube /
