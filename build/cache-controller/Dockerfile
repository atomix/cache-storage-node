FROM alpine:3.9

USER nobody

ADD build/cache-controller/_output/cache-controller /usr/local/bin/cache-controller

ENTRYPOINT ["cache-controller"]
