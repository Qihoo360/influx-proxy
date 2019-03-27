### Makefile ---

## Author: Shell.Xu
## Version: $Id: Makefile,v 0.0 2017/01/17 03:44:24 shell Exp $
## Copyright: 2017, Eleme <zhixiang.xu@ele.me>
## License: MIT
## Keywords:
## X-URL:

# Build variables
REGISTRY_URI :=wilhelmguo
RELEASE_VERSION :=$(shell git describe --always --tags)

all: build-image push-image

build:
	mkdir -p bin
	go build -o bin/influx-proxy github.com/wilhelmguo/influx-proxy/service

test:
	go test -v github.com/wilhelmguo/influx-proxy/backend

bench:
	go test -bench=. github.com/wilhelmguo/influx-proxy/backend

clean:
	rm -rf bin

build-image:
	@echo "version: $(RELEASE_VERSION)"
	docker build --no-cache -t $(REGISTRY_URI)/influx-proxy:$(RELEASE_VERSION) .

push-image:
	docker push $(REGISTRY_URI)/influx-proxy:$(RELEASE_VERSION)

### Makefile ends here
