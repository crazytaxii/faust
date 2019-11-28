PHONY: build

TARGET ?= faust
ODIR ?= .

build:
	go build -ldflags="-X 'main.AppVersion=$(shell git rev-parse --short HEAD)'" -o $(ODIR)/$(TARGET) .
