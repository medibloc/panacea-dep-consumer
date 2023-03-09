export GO111MODULE = on

build_tags := $(strip $(BUILD_TAGS))
BUILD_FLAGS := -tags "$(build_tags)"

OUT_DIR = ./build

.PHONY: build test install clean

build: go.sum
	go build -mod=readonly $(BUILD_FLAGS) -o $(OUT_DIR)/consumerd ./cmd/consumerd

test:
	go test -v ./...

install: go.sum
	go install -mod=readonly $(BUILD_FLAGS) ./cmd/consumerd

clean:
	go clean
	rm -rf $(OUT_DIR)