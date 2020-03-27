PREFIX?=./build

install:
	go install ./...

.PHONY: build
build:
	mkdir -p $(PREFIX)
	env go build -o "$(PREFIX)/smt.logger" ./cmd/smt.logger

.PHONY: build.linux.nat
build.linux:
	mkdir -p $(PREFIX)
	env GOOS=linux go build -o "$(PREFIX)/smt.logger.linux" ./cmd/smt.logger

.PHONY: build.linux.arm
build.linux.arm:
	mkdir -p $(PREFIX)
	env GOOS=linux GOARCH=arm go build -o "$(PREFIX)/smt.logger.linux.arm" ./cmd/smt.logger
