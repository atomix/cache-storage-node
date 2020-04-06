FROM alpine:3.9

USER nobody

ADD build/cache-storage-controller/_output/local-replica /usr/local/bin/cache-storage-controller

ENTRYPOINT ["cache-storage-controller"]
