# MIT License
# Copyright (c) 2025 Toni Liesche
#
# Permission is hereby granted, free of charge, to any person obtaining a copy
# of this software and associated documentation files (the "Software"), to deal
# in the Software without restriction, including without limitation the rights
# to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
# copies of the Software, and to permit persons to whom the Software is
# furnished to do so, subject to the following conditions:
#
# The above copyright notice and this permission notice shall be included in all
# copies or substantial portions of the Software.

FROM golang:1.24-bookworm AS build-stage

ARG CACHEBUST
RUN echo ${CACHEBUST} && \
    git clone https://github.com/toniliesche/powerdns-auth-proxy.git /tmp/powerdns-auth-proxy

WORKDIR /tmp/powerdns-auth-proxy

ARG BUILDVER
RUN --mount=type=cache,target=/root/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    git checkout ${BUILDVER} && \
    make build

FROM debian:bookworm-slim

ARG BUILDVER
LABEL org.opencontainers.image.authors="Toni Liesche <docker@toniliesche.de>"
LABEL org.opencontainers.image.url="https://hub.docker.com/r/tliesche/powerdns-auth-proxy"
LABEL org.opencontainers.image.documentation="https://raw.githubusercontent.com/toniliesche/powerdns-auth-proxy/refs/heads/develop/README.md"
LABEL org.opencontainers.image.source="https://github.com/toniliesche/powerdns-auth-proxy"
LABEL org.opencontainers.image.version=${BUILDVER}
LABEL org.opencontainers.image.licenses="MIT"
LABEL org.opencontainers.image.title="PowerDNS Authentication Proxy"
LABEL org.opencontainers.image.description="Database backed authentication proxy for PowerDNS HTTP API"

ENV GIN_MODE=release

RUN apt update && apt install -y --no-install-recommends ca-certificates procps htop
RUN adduser --disabled-password --shell /usr/sbin/nologin --home /var/lib/powerdns-auth-proxy apiuser

COPY --from=build-stage /tmp/powerdns-auth-proxy/gateway /usr/local/bin/powerdns-auth-proxy
COPY --from=build-stage /tmp/powerdns-auth-proxy/docker/auth-proxy/files/entry-point.sh /entry-point.sh
COPY --from=build-stage /tmp/powerdns-auth-proxy/docker/auth-proxy/files/configure.d /configure.d
COPY --from=build-stage /tmp/powerdns-auth-proxy/docker/auth-proxy/files/var/lib/powerdns-auth-proxy/res /var/lib/powerdns-auth-proxy/res

ENTRYPOINT [ "/entry-point.sh" ]
CMD [ "su", "apiuser", "-s", "/usr/bin/bash", "-c", "powerdns-auth-proxy --config /etc/powerdns-auth-proxy/config.yaml run" ]
