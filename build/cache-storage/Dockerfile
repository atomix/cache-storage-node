FROM alpine:3.9

USER nobody

ADD build/cache-storage/_output/cache-storage /usr/local/bin/cache-storage

ENTRYPOINT ["cache-storage"]
