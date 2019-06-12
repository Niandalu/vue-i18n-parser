.PHONY: test
test:
	go test ./...

example-collect:
	go run cmd/vip.go c --diff ./testdata

example-feed:
	go run cmd/vip.go f ./testdata/translated.csv

build:
	dep ensure
	GOOS=linux GOARCH=amd64 go build cmd/vip.go
