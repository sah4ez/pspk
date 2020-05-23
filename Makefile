VERSION=0.1.16
NAME=pspk
GIT_REV?=$(shell git rev-parse --short HEAD)
LDFLAGS=-ldflags "-X main.Version=$(VERSION) -X main.Hash=$(GIT_REV) -X main.BuildDate=$(BUILD_DATE)"
GO=GO111MOUDLE=on go
SIGNATORY=pspk-sign


.PHONY: build
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

wasm_exec.js: ./web_wasm/wasm_exec.js
	cp "${GOROOT}/misc/wasm/wasm_exec.js" ./web_wasm/

.PHONY: wasm
wasm: wasm_exec.js
	env GOOS=js GOARCH=wasm go build ${LDFLAGS} -o ./bin/wasm/${NAME}.wasm ./cmd/wasm/
	env GOOS=js GOARCH=wasm go build ${LDFLAGS} -o ./bin/wasm/publish ./cmd/wasm/publish/.

local-web: wasm
	cp ./bin/wasm/* ./web_wasm/

.PHONY: clean
clean:
	rm -rf _build

.PHONY: release_web
release_web:
	scp -r ./web/*  ${FREECONTENT_SPACE}:/var/www/pspk/
