DB_SOURCE=root:example@tcp(127.0.0.1:3306)
APP_NAME=playground

.PHONY: gin

gin:
	go build -o bin/gin ./cmd/gin/main.go

migrateup:
	migrate -path db/migration -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose up

migratedown:
	migrate -path db/migration -database 'mysql://$(DB_SOURCE)/$(APP_NAME)' -verbose down
