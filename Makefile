DB_SOURCE=root:example@tcp(127.0.0.1:3306)
APP_NAME=playground

.PHONY: gin migrateup migratedown sqlc

gin:
	go build -o bin/gin ./cmd/gin/main.go

migrateup:
	migrate -path tools/migration -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose up

migratedown:
	migrate -path tools/migration -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose down

sqlc:
	sqlc generate -f ./sqlc.yaml
