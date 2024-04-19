TARGET ?= faust
ODIR ?= dist

PHONY: build install

build:
	go build -ldflags="-X 'main.AppVersion=$(shell git rev-parse --short HEAD)'" -o $(ODIR)/$(TARGET) ./cmd/faust

install: build
	cp $(ODIR)/$(TARGET) /usr/local/bin/
