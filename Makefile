VERSION := 0.1.0
BUILD_DATE := $(shell date +%F)

all: fusili

fusili:
	go build -mod=vendor -ldflags " \
		-X main.version=$(VERSION) \
		-X main.buildDate=$(BUILD_DATE) \
		" ./cmd/fusili
