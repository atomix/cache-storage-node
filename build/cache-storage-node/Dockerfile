FROM alpine:3.9

USER nobody

ADD build/cache-storage-node/_output/cache-storage-node /usr/local/bin/cache-storage-node

ENTRYPOINT ["cache-storage-node"]
