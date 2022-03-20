build: test
	go build
test:
	go test
start: build
	./gannant-code-test
all: build test start
