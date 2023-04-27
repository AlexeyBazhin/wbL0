GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)
OUT := bin/wbl0
OUT_PATH=$(shell pwd)/bin/$(GOOS)_$(GOARCH)

clean:
	rm -rf ./bin/*
.PHONY: clean

clean.bin: ## remove $(OUT_PATH) directory
	rm -rf $(OUT_PATH)
.PHONY: clean.bin

build: clean
	GOOS=linux GOARCH=amd64 go build -ldflags "-w -s" -o $(OUT) ./cmd/wbl0
.PHONY: build

dbuild:
	docker compose -f ../wbL0/docker-compose.yml build --no-cache
.PHONY: dbuild

dstart:
	docker compose -f ../wbL0/docker-compose.yml up --force-recreate
.PHONY: dstart

run: dbuild dstart
.PHONY: run