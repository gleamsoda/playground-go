# Build stage
FROM golang:1.21-bullseye AS builder
WORKDIR /app
# Copy the whole directory to the WORKDIR
COPY . .
RUN CGO_ENABLED=0 go build -o ./bin/playground ./cmd/playground/main.go

# Run stage
FROM alpine:3.18
WORKDIR /app
# Copy the binary to the WORKDIR
COPY --from=builder /app/bin/playground ./playground
COPY wait-for.sh .
COPY db/migrations ./db/migrations

EXPOSE 8080
CMD [ "gin" ]
ENTRYPOINT [ "/app/playground" ]