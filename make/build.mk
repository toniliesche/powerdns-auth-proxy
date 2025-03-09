.PHONY: build
build:
	go build -o build/powerdns-auth-proxy main/main.go

build-docker-rc: set-version-rc
	docker buildx build \
		--pull \
		--platform linux/amd64,linux/arm64 \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=$(run.build.version) \
		-t tliesche/powerdns-auth-proxy:$(run.build.version) \
		$(if $PUSH,--push,) \
 		docker/auth-proxy

build-docker-%: set-version-%
	docker buildx build \
		--pull \
		--platform linux/amd64,linux/arm64 \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=$(run.build.version) \
		-t tliesche/$*:$(run.build.version) \
		-t tliesche/$*:$(run.build.version.minor) \
		-t tliesche/$*:$(run.build.version.major) \
		-t tliesche/$*:latest \
		$(if $PUSH,--push,) \
		docker/auth-proxy

build-dev-docker:
	docker build \
		--pull \
		--build-arg CACHEBUST=$(shell date +%s) \
		--build-arg BUILDVER=develop \
		-t tliesche/powerdns-auth-proxy:develop \
		docker/auth-proxy
