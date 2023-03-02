export GO111MODULE = on

build_tags := $(strip $(BUILD_TAGS))
BUILD_FLAGS := -tags "$(build_tags)"

OUT_DIR = ./build

MODULE=github.com/medibloc/panacea-dep-consumer

.PHONY: build test clean

build: go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/dep-consumer

test:
	go test -v ./...

clean:
	go clean
	rm -rf $(OUT_DIR)