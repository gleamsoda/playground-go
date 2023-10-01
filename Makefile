APP_NAME=playground
GRPC_BASE=internal/delivery/grpc
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: build test cover mock sqlc proto migrate/up migrate/down evans

build:
	go build -o bin/playground ./cmd/playground/main.go

test:
	go test -v -cover ./... -coverprofile=$(COVERAGE_OUT)

cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

mock:
	mockgen -source ./internal/wallet/repository.go -destination ./test/mock/wallet/repository.go
	mockgen -source ./internal/wallet/usecase.go -destination ./test/mock/wallet/usecase.go
	mockgen -source ./internal/wallet/mq/mq.go -destination ./test/mock/wallet/mq/mq.go
	mockgen -source ./internal/pkg/token/manager.go -destination ./test/mock/pkg/token/manager.go

sqlc:
	sqlc generate -f ./sqlc.yaml

proto:
	rm -f $(GRPC_BASE)/gen/*.pb.go $(GRPC_BASE)/gen/*.pb.gw.go
	rm -f docs/*.swagger.json
	protoc --proto_path=$(GRPC_BASE)/proto \
		--go_out=$(GRPC_BASE)/gen --go_opt=paths=source_relative \
		--go-grpc_out=$(GRPC_BASE)/gen --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GRPC_BASE)/gen --grpc-gateway_opt paths=source_relative \
		--openapiv2_out=./docs --openapiv2_opt=allow_merge=true,merge_file_name=$(APP_NAME) \
		$(GRPC_BASE)/proto/*.proto

migrate/up:
	migrate -path db/migrations -database 'mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/$(APP_NAME)' -verbose up

migrate/down:
	migrate -path db/migrations -database 'mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/$(APP_NAME)' -verbose down

evans:
	evans --host localhost --port 9090 -r repl
