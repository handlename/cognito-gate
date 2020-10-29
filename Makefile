VERSION=$(shell cat ./VERSION)

cmd/cognito-gate/cognito-gate: go.sum *.go */**/*.go
	go build -v -ldflags '-X main.version=$(VERSION)' -o $@ cmd/cognito-gate/main.go

test:
	go test -v ./...
