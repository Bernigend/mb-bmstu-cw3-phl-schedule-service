# Название сервиса
SERVICE_NAME = schedule-service
# Текущая директория
PWD = $(shell pwd)

# Директория с .proto файлами
PROTO_FILES_DIR = api
# Директория для .pb.go файлов (после компиляции .proto)
GO_OUT_DIR = pkg/$(SERVICE_NAME)-api
# Название .proto файла для компиляции
PROTO_FILES = $(SERVICE_NAME)-api.proto

# 8 символов последнего коммита
LAST_COMMIT_HASH = $(shell git rev-parse HEAD | cut -c -8)

# Компиляция proto файлов
.PHONY: generate
generate:
	mkdir -p $(PWD)/pkg/$(SERVICE_NAME)-api && \
	cd $(PROTO_FILES_DIR) && \
	protoc -I. --go_out=plugins=grpc:$(PWD)/$(GO_OUT_DIR) $(PROTO_FILES) && \
	echo "New pb files generated"