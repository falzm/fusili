VERSION := 0.1.0
BUILD_DATE := $(shell date +%F)

all: fusili

fusili:
	gb build -ldflags " \
		-X main.version=$(VERSION) \
		-X main.buildDate=$(BUILD_DATE) \
		" fusili/fusili

clean:
	rm -rf bin/ pkg/
