VERSION=$(shell cat ./VERSION)
PROJECT_USERNAME=handlename
PROJECT_REPONAME=cognito-gate
DIST_DIR=dist

cmd/cognito-gate/cognito-gate: go.sum *.go */**/*.go
	go build -v -ldflags '-X main.version=$(VERSION)' -o $@ cmd/cognito-gate/main.go

test:
	go test -v ./...

.PHONY: release
release:
	-git tag v$(VERSION)
	git push
	git push --tags

.PHONY: dist
dist: clean
	CGO_ENABLED=0 goxz \
	  -pv 'v$(VERSION)' \
	  -n cognito-gate \
	  -build-ldflags '-X main.version=$(VERSION)' \
	  -os='linux' \
	  -arch='amd64' \
	  -d $(DIST_DIR) \
	  ./cmd/cognito-gate

.PHONY: upload
upload: dist
	mkdir -p $(DIST_DIR)
	ghr \
	  -u '$(PROJECT_USERNAME)' \
	  -r '$(PROJECT_REPONAME)' \
	  -prerelease \
	  -replace \
	  'v$(VERSION)' \
	  $(DIST_DIR)

clean:
	rm -rf cmd/cognito-gate/cognito-gate $(DIST_DIR)/*
