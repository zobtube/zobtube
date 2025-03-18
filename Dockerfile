FROM busybox AS build-env

RUN mkdir -p /build/tmp
RUN chmod 777 /build/tmp

FROM scratch

ENTRYPOINT ["/zobtube"]

COPY --from=build-env /build/tmp /tmp
COPY zobtube /
