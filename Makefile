ifeq ($(VERSION),)
  VERSION_TAG=$(shell git describe --abbrev=0 --tags --exact-match 2>/dev/null || echo latest)
else
  VERSION_TAG=$(VERSION)
endif


.PHONY: build
build:
	go build -o ./.bin/tenant-controller -ldflags "-X \"main.version=$(VERSION_TAG)\""  main.go
