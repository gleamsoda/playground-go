APP_NAME=playground
GRPC_BASE=internal/delivery/grpc
COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html

.PHONY: build test cover generate sqlc proto evans migrate-up migrate-down migrate-create

build:
	go build -o bin/playground ./cmd/playground/main.go

test:
	go test -v -short -cover ./... -coverprofile=$(COVERAGE_OUT)

cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

generate:
	go generate ./...

sqlc:
	rm -f ./internal/app/repository/sqlc/gen/*.go
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

evans:
	evans --host localhost --port 9090 -r repl

migrate-up:
	migrate -path db/migrations -database 'mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/$(APP_NAME)' -verbose up

migrate-down:
	migrate -path db/migrations -database 'mysql://${DB_USER}:${DB_PASSWORD}@tcp(${DB_HOST}:${DB_PORT})/$(APP_NAME)' -verbose down

migrate-create:
	migrate create -ext sql -dir db/migrations -seq $(name)
