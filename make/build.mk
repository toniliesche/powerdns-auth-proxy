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

GOLANGVER=1.24.1

.PHONY: build
build:
	go build -o build/powerdns-auth-proxy main/main.go

build-docker-rc: set-version-rc set-commit
	docker buildx build \
		--pull \
		--platform linux/amd64,linux/arm64 \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=$(run.build.version) \
		--build-arg COMMIT=$(run.commit) \
		--build-arg GOLANGVER=$(GOLANGVER) \
		-t tliesche/powerdns-auth-proxy:edge \
		-t tliesche/powerdns-auth-proxy:$(run.build.version) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.rc) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.rc.minor) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.rc.major) \
		-t tliesche/powerdns-auth-proxy:$(run.commit) \
		$(if $(PUSH),--push,--no-cache --progress=plain) \
		docker/auth-proxy

build-docker-%: set-version-% set-commit
	docker buildx build \
		--pull \
		--platform linux/amd64,linux/arm64 \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=$(run.build.version) \
		--build-arg COMMIT=$(run.commit) \
		--build-arg GOLANGVER=$(GOLANGVER) \
		-t tliesche/powerdns-auth-proxy:latest \
		-t tliesche/powerdns-auth-proxy:$(run.build.version) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.minor) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.major) \
		-t tliesche/powerdns-auth-proxy:$(run.commit) \
		$(if $(PUSH),--push,--no-cache --progress=plain) \
		docker/auth-proxy

build-dev-docker: set-commit
	docker build \
		--pull \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=develop \
		--build-arg COMMIT=$(run.commit) \
		--build-arg GOLANGVER=$(GOLANGVER) \
		-t tliesche/powerdns-auth-proxy:develop \
		docker/auth-proxy
