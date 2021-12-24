REPO = ghcr.io/jaeg/shorten
BINARY = shorten

TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell git log -1 --format=%cd --date=format:"%Y%m%d")
VERSION := $(TAG:v%=%)
ifneq ($(COMMIT), $(TAG_COMMIT))
	VERSION := $(VERSION)-$(COMMIT)-$(DATE)
endif
ifeq ($(VERSION),)
	VERSION := $(COMMIT)-$(DATA)
endif

ifneq ($(shell git status --porcelain),)
	VERSION := $(VERSION)-dirty
endif

bin:
	mkdir bin

vendor:
	go mod vendor

image: build-linux
	docker build -t $(REPO):$(VERSION) . --build-arg binary=$(BINARY)-linux --build-arg version=$(VERSION)

image-pi: build-linux-pi

	docker build -t $(REPO):$(VERSION)-pi . --build-arg binary=$(BINARY)-linux-pi --build-arg version=$(VERSION)

run:
	go run -mod=vendor .

build: bin
	go build -mod=vendor -o ./bin/$(BINARY)

build-linux: bin
	env GOOS=linux GOARCH=amd64 go build -mod=vendor -o ./bin/$(BINARY)-linux

build-linux-pi: bin
	env GOOS=linux GOARCH=arm GOARM=7 go build -mod=vendor -o ./bin/$(BINARY)-linux-pi

publish-pi:
	docker push $(REPO):$(VERSION)-pi
	docker tag $(REPO):$(VERSION)-pi $(REPO):latest-pi
	docker push $(REPO):latest-pi

publish:
	docker push $(REPO):$(VERSION)
	docker tag $(REPO):$(VERSION) $(REPO):latest
	docker push $(REPO):latest

.PHONY: update-go-deps
update-go-deps:
	@echo ">> updating Go dependencies"
	@for m in $$(go list -mod=readonly -m -f '{{ if and (not .Indirect) (not .Main)}}{{.Path}}{{end}}' all); do \
		go get $$m; \
	done
	go mod tidy
ifneq (,$(wildcard vendor))
	go mod vendor
endif