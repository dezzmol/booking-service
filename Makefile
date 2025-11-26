PROTO_DIR := ./api
EXTERNAL_PROTO_DIR := ./api/external
GEN_DIR := ./internal/generated
MODULE_NAME := github.com/dezzmol/booking-service/internal/generated

PROTO_INCLUDES := -I$(PROTO_DIR) -I./vendor

.PHONY: gen-proto
gen-proto:
	@echo "Generating proto stubs..."
	@mkdir -p $(GEN_DIR)
	@protoc $(PROTO_INCLUDES) \
		--go_out=$(GEN_DIR) --go_opt=module=$(MODULE_NAME) \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=module=$(MODULE_NAME) \
		--grpc-gateway_out=$(GEN_DIR) --grpc-gateway_opt=module=$(MODULE_NAME) \
		--grpc-gateway_opt=allow_delete_body=true \
		--openapiv2_out=$(GEN_DIR) \
		$(shell find $(PROTO_DIR) -name '*.proto')

	@protoc $(PROTO_INCLUDES) \
		--go_out=$(GEN_DIR) --go_opt=module=$(MODULE_NAME) \
		--go-grpc_out=$(GEN_DIR) --go-grpc_opt=module=$(MODULE_NAME) \
		--grpc-gateway_out=$(GEN_DIR) --grpc-gateway_opt=module=$(MODULE_NAME) \
		--grpc-gateway_opt=allow_delete_body=true \

		$(shell find $(EXTERNAL_PROTO_DIR) -name '*.proto')
	@echo "Proto stubs generated in $(GEN_DIR)"

.PHONY: clean-proto
clean-proto:
	@echo "Cleaning generated proto stubs..."
	@rm -rf $(GEN_DIR)/*
	@echo "Clean complete"

.PHONY: install-deps
install-deps:
	@echo "Installing protoc dependencies..."
	@go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	@go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	@echo "Dependencies installed"

.PHONE: build
build:
	@echo "Building service"
	@go build -v ./cmd/main.go
