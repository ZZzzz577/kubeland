FROM docker.io/library/golang:1.24.1 AS builder

ENV GOPROXY=https://goproxy.cn

COPY . /src
WORKDIR /src

RUN go mod download &&  \
    mkdir -p bin && \
    go build -o ./bin/server ./...

FROM docker.io/library/debian:stable-slim

RUN apt-get update && apt-get install -y --no-install-recommends \
		ca-certificates \
        netbase \
        netavark \
        buildah \
#        fuse-overlayfs \
        && rm -rf /var/lib/apt/lists/ \
        && apt-get autoremove -y && apt-get autoclean -y \
        && echo 'unqualified-search-registries = ["docker.io"]' > /etc/containers/registries.conf \
        && echo '[storage]' > /etc/containers/storage.conf \
        && echo 'driver = "vfs"' >> /etc/containers/storage.conf \
        && echo 'runroot = "/run/containers/storage"' >> /etc/containers/storage.conf \
        && echo 'graphroot = "/var/lib/containers/storage"' >> /etc/containers/storage.conf 
#        && echo '[storage.options.overlay]' >> /etc/containers/storage.conf \
#        && echo 'mount_program = "/usr/bin/fuse-overlayfs"' >> /etc/containers/storage.conf

COPY --from=builder /src/bin /app

WORKDIR /app

CMD ["./server"]