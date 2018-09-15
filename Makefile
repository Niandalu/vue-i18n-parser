.PHONY: test
test:
	go test ./...

example-collect:
	go run cmd/vip.go c --diff ./testdata

example-feed:
	go run cmd/vip.go f ./testdata/translated.csv

build:
	dep ensure
	go build cmd/vip.go
