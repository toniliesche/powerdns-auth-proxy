test:
	@go test ./...

test-verbose:
	@go test -v ./...

test-coverage:
	@go test -coverprofile=coverage.out.tmp ./...
	@cat coverage.out.tmp | grep -v "mocks" > coverage.out.tmp.2
	@go tool cover -html coverage.out.tmp.2 -o cover.html
	@rm -rf coverage.out*
