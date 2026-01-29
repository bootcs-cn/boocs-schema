.PHONY: build install test clean release

BINARY_NAME=bootcs-validate
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

build:
	go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY_NAME) ./cmd/bootcs-validate

install:
	go install -ldflags "-X main.version=$(VERSION)" ./cmd/bootcs-validate

test:
	go test -v ./...

clean:
	rm -rf bin/

# 验证 bcs100x 课程（测试用）
validate-bcs100x:
	go run ./cmd/bootcs-validate ../bootcs-courses/bcs100x

# 构建多平台
release:
	GOOS=darwin GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY_NAME)-darwin-amd64 ./cmd/bootcs-validate
	GOOS=darwin GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY_NAME)-darwin-arm64 ./cmd/bootcs-validate
	GOOS=linux GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY_NAME)-linux-amd64 ./cmd/bootcs-validate
	GOOS=linux GOARCH=arm64 go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY_NAME)-linux-arm64 ./cmd/bootcs-validate
	GOOS=windows GOARCH=amd64 go build -ldflags "-X main.version=$(VERSION)" -o bin/$(BINARY_NAME)-windows-amd64.exe ./cmd/bootcs-validate
