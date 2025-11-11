# Makefile
include vendor.proto.mk

# Используем bin в текущей директории для установки плагинов protoc
LOCAL_BIN := $(CURDIR)/bin
BUF_BUILD := $(LOCAL_BIN)/buf

# устанавливаем необходимые плагины
.bin-deps: export GOBIN := $(LOCAL_BIN)
.bin-deps:
	@echo "Installing binary dependencies..."
	@go install github.com/bufbuild/buf/cmd/buf@v1.41.0
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.buf-generate:
	PATH="$(LOCAL_BIN):$(PATH)" $(BUF_BUILD) generate

generate: .buf-generate .tidy

.tidy:
	go mod tidy
# Объявляем, что текущие команды не являются файлами и инструментируем Makefile не искать изменения в файловой системе
.PHONY: .bin-deps