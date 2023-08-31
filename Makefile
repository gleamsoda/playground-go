COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html
DB_SOURCE=root:example@tcp(127.0.0.1:3306)
APP_NAME=playground
GRPC_BASE=driver/grpc

.PHONY: build test cover mock sqlc proto migrate/up migrate/down

build:
	go build -o bin/playground ./cmd/playground/main.go

test:
	go test -v -cover ./... -coverprofile=$(COVERAGE_OUT)

cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

mock:
	mockgen playground/app Repository > ./test/mock/app/repository.go
	mockgen playground/app Usecase > ./test/mock/app/usecase.go
	mockgen playground/pkg/token Manager > ./test/mock/token/manager.go

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
	migrate -path tools/migrations -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose up

migrate/down:
	migrate -path tools/migrations -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose down
