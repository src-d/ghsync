# Package configuration
PROJECT = ghsync
COMMANDS = cmd/ghsync

GO_BUILD_ENV = CGO_ENABLED=0
PKG_OS = darwin linux

# Add prerequisites
build: vendor bindata

# Including ci Makefile
CI_REPOSITORY ?= https://github.com/src-d/ci.git
CI_BRANCH ?= v1
CI_PATH ?= .ci
MAKEFILE := $(CI_PATH)/Makefile.main
$(MAKEFILE):
	git clone --quiet --depth 1 -b $(CI_BRANCH) $(CI_REPOSITORY) $(CI_PATH);
-include $(MAKEFILE)

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: bindata
bindata:
	go get -u github.com/jteeuwen/go-bindata/... && \
	go-bindata \
		-pkg migrations \
		-o ./models/migrations/bindata.go \
		-prefix models/sql \
		-modtime 1 \
		models/sql/...
