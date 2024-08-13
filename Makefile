PROTO_DIR=test/proto
PB_OUT_DIR=pb
PROTOC_VERSION=21.12

.PHONY: all proto clean install-protoc

all: proto

# Generate Go code from proto files
proto:
	@mkdir -p $(PB_OUT_DIR)
	protoc --go_out=$(PB_OUT_DIR) --go_opt=paths=source_relative $(PROTO_DIR)/*.proto

# Install protoc
install-protoc:
	# Linux
#	@curl -OL https://github.com/protocolbuffers/protobuf/releases/download/v$(PROTOC_VERSION)/protoc-$(PROTOC_VERSION)-linux-x86_64.zip
#	@unzip -o protoc-$(PROTOC_VERSION)-linux-x86_64.zip -d /usr/local protoc && rm protoc-$(PROTOC_VERSION)-linux-x86_64.zip

	# MacOS
	@brew install protobuf

# クリーンアップ
clean:
	rm -rf $(PB_OUT_DIR)/*.go