update-deps:
	@go get -u ./...

get-bash:
	@$(AUTH_PROXY_CMD) exec api bash

logs:
	@$(AUTH_PROXY_CMD) logs api -f --no-log-prefix
