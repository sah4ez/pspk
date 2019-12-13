VERSION=0.1.8
NAME=pspk
GIT_REV?=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Hash=$(GIT_REV) -X main.BuildDate=$(BUILD_DATE)"
GO=GO111MOUDLE=on go


build:
	VERSION=$(VERSION)-dev $(GO) build $(LDFLAGS) -o bin/${NAME} ./cmd/cli/

release: clean
	mkdir -p _build
	GOOS=darwin  GOARCH=amd64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-$(VERSION)-darwin-amd64 ./cmd/cli
	GOOS=linux   GOARCH=amd64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-$(VERSION)-linux-amd64 ./cmd/cli
	GOOS=linux   GOARCH=arm   $(GO) build $(LDFLAGS) -o _build/$(NAME)-$(VERSION)-linux-arm ./cmd/cli
	GOOS=linux   GOARCH=arm64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-$(VERSION)-linux-arm64 ./cmd/cli
	GOOS=darwin  GOARCH=amd64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-srv-$(VERSION)-darwin-amd64 ./cmd/server
	GOOS=linux   GOARCH=amd64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-srv-$(VERSION)-linux-amd64 ./cmd/server
	GOOS=linux   GOARCH=arm   $(GO) build $(LDFLAGS) -o _build/$(NAME)-srv-$(VERSION)-linux-arm ./cmd/server
	GOOS=linux   GOARCH=arm64 $(GO) build $(LDFLAGS) -o _build/$(NAME)-srv-$(VERSION)-linux-arm64 ./cmd/server
	cd _build; sha256sum * > sha256sums.txt

clean:
	rm -rf _build

release_web:
	scp ./web/*  ${FREECONTENT_SPACE}:/var/www/pspk/
