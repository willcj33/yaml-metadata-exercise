.PHONY: build start test vet init ensure

init:
	dep init

ensure:
	dep ensure -update

vet:
	go vet -v ./...

build: vet
	go build

start: vet build
	./yaml-metadata-exercise

test: vet build
	go test -v ./...
