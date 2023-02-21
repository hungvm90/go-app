.DEFAULT_GOAL := build
fmt:
	go fmt ./...
.PHONY:fmt
lint: fmt
	golint ./...
.PHONY:lint
vet: fmt
	go vet ./...
.PHONY:vet
build: fmt
	go test ./... -coverprofile=coverage.out
.PHONY:build