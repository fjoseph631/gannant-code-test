build:
	go build
test: build
	go test
start: build
	./gannant-code-test
all: build test start
