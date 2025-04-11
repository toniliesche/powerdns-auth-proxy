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

test:
	@go test ./...

test-verbose:
	@go test -v ./...

test-coverage:
	@go test -coverprofile=coverage.out.tmp ./...
	@cat coverage.out.tmp | grep -v "mocks" > coverage.out.tmp.2
	@go tool cover -html coverage.out.tmp.2 -o cover.html
	@rm -rf coverage.out*
