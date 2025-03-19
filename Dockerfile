FROM busybox AS build-env

RUN mkdir -p /build/tmp

FROM scratch

ENTRYPOINT ["/zobtube"]

COPY --from=build-env --chmod 777 /build/tmp /tmp
COPY zobtube /
