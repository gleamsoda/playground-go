package server

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"playground/app/repository"
	"playground/app/usecase"
	"playground/cmd/grpc/server/gen"
	"playground/config"
	"playground/pkg/token"
)

// NewServer creates a new gRPC server.
func NewServer(cfg config.Config) (*grpc.Server, error) {
	conn, err := sql.Open("mysql", cfg.DBSource)
	if err != nil {
		return nil, err
	}
	if err := conn.Ping(); err != nil {
		return nil, err
	}
	tm, err := token.NewPasetoMaker(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}

	userRepo := repository.NewUserRepository(conn)
	sessionRepo := repository.NewSessionRepository(conn)
	userUsecase := usecase.NewUserUsecase(userRepo, sessionRepo, tm, cfg)

	ctrl := NewController(userUsecase)
	svr := grpc.NewServer()
	gen.RegisterPlaygroundServer(svr, ctrl)
	reflection.Register(svr)

	return svr, nil
}
