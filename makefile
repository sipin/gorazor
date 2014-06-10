export GOPATH=$(shell pwd)

test:
	go get github.com/sipin/gorazor
	go test -v ./gorazor
