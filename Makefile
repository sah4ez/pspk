VERSION=0.0.1
NAME=pspk
GIT_REV?=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Hash=$(GIT_REV) -X main.BuildDate=$(BUILD_DATE)"
GO=GO111MOUDLE=on go


build:
	$(GO) build -o bin/${NAME} ./cmd/cli/main.go

release: clean
	mkdir -p _build
	GOOS=darwin  GOARCH=amd64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-$(VERSION)-darwin-amd64 ./cmd/cli
	GOOS=linux   GOARCH=amd64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-$(VERSION)-linux-amd64 ./cmd/cli
	GOOS=linux   GOARCH=arm   $(GO) build $(LDFLAGS) -o _build/$(NAME)-$(VERSION)-linux-arm ./cmd/cli
	GOOS=linux   GOARCH=arm64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-$(VERSION)-linux-arm64 ./cmd/cli
	cd _build; sha256sum * > sha256sums.txt

clean:
	rm -rf _build
