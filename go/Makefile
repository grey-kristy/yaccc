
all: build

fmt:
	gofmt -w -s ./*.go
	gofmt -w -s ./convertor/*.go

build:
	go build -o ./bin/yaccc yaccc.go
