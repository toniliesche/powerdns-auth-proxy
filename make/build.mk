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
		-t tliesche/powerdns-auth-proxy:$(run.commit) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.rc) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.rc.minor) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.rc.major) \
		-t tliesche/powerdns-auth-proxy:edge \
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
		-t tliesche/powerdns-auth-proxy:$(run.commit) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.minor) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version.major) \
		-t tliesche/powerdns-auth-proxy:latest \
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
