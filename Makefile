.PHONY: install test-dev test cover run-dev build

install:
	go mod download

test-dev:
	mkdir -p ./coverage && \
		go test -v -coverprofile=./coverage/coverage.out -covermode=atomic ./...

test:
	mkdir -p ./coverage && \
		go test -v -coverprofile=./coverage/coverage.out -covermode=atomic ./...

cover: test-dev
	go tool cover -func=./coverage/coverage.out &&\
		go tool cover -html=./coverage/coverage.out -o ./coverage/coverage.html

run-dev:
	go run ./main.go

build:
	CGO_ENABLED=0 GOOS=linux go build -a -o app &&\
		cp app /tmp/app