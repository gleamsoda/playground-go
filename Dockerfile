# Build stage
FROM golang:1.21-bullseye AS builder
WORKDIR /app
# Copy the whole directory to the WORKDIR
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/gin ./cmd/gin/main.go
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.18
WORKDIR /app
# Copy the binary to the WORKDIR
COPY --from=builder /app/bin/gin ./gin
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY wait-for .
COPY tools/migrations ./migrations

EXPOSE 8080
CMD [ "gin" ]
ENTRYPOINT [ "/app/start.sh", "/app/playground" ]