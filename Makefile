COVERAGE_OUT=coverage.out
COVERAGE_HTML=coverage.html
DB_SOURCE=root:example@tcp(127.0.0.1:3306)
APP_NAME=playground
GRPC_BASE=cmd/grpc/internal

.PHONY: gin grpc test cover migrateup migratedown sqlc proto

gin:
	go build -o bin/gin ./cmd/gin/main.go

grpc:
	go build -o bin/grpc ./cmd/grpc/main.go

test:
	go test -v -cover ./... -coverprofile=$(COVERAGE_OUT)

cover:
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)

migrateup:
	migrate -path tools/migration -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose up

migratedown:
	migrate -path tools/migration -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose down

sqlc:
	sqlc generate -f ./sqlc.yaml

proto:
	rm -f $(GRPC_BASE)/boundary/*.pb.go $(GRPC_BASE)/boundary/*.pb.gw.go
	rm -f doc/*.swagger.json
	protoc --proto_path=$(GRPC_BASE)/proto \
		--go_out=$(GRPC_BASE)/boundary --go_opt=paths=source_relative \
    	--go-grpc_out=$(GRPC_BASE)/boundary --go-grpc_opt=paths=source_relative \
		--grpc-gateway_out=$(GRPC_BASE)/boundary --grpc-gateway_opt paths=source_relative \
		--openapiv2_out=./doc --openapiv2_opt=allow_merge=true,merge_file_name=$(APP_NAME) \
    	$(GRPC_BASE)/proto/*.proto
