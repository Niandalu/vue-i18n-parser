.PHONY: test
test:
	go test ./...

example-collect:
	go run cmd/vip.go c --diff ./testdata

example-feed:
	go run cmd/vip.go f ./testdata/translated.csv

build:
	GOOS=linux GOARCH=amd64 go build -o vip.linux cmd/vip.go
	GOOS=darwin GOARCH=amd64 go build -o vip.darwin cmd/vip.go
