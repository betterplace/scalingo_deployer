DOCKER_IMAGE = scalingo_deployer
PROJECT_ID = betterplace-183212
REMOTE_TAG = eu.gcr.io/$(PROJECT_ID)/$(DOCKER_IMAGE)
GOPATH := $(shell pwd)/gospace
GOBIN = $(GOPATH)/bin

.EXPORT_ALL_VARIABLES:

all: scalingo_deployer

scalingo_deployer: cmd/scalingo_deployer/main.go *.go
	go build -o $@ $<

setup: fake-package
	go mod download

fake-package:
	rm -rf $(GOPATH)/src/github.com/betterplace/scalingo_deployer
	mkdir -p $(GOPATH)/src/github.com/betterplace
	ln -s $(shell pwd) $(GOPATH)/src/github.com/betterplace/scalingo_deployer

clean:
	@rm -f scalingo_deployer tags

clobber: clean
	@rm -rf $(GOPATH)/*

tags: clean
	@gotags -tag-relative=false -silent=true -R=true -f $@ . $(GOPATH)

build-info:
	@echo $(REMOTE_TAG)

build:
	docker build -t $(DOCKER_IMAGE) .
	$(MAKE) build-info

build-force:
	docker build -t $(DOCKER_IMAGE) --no-cache .
	$(MAKE) build-info

debug:
	docker run --rm -it $(DOCKER_IMAGE) bash

push: build
	gcloud auth configure-docker
	docker tag $(DOCKER_IMAGE) $(REMOTE_TAG)
	docker push $(REMOTE_TAG)
