package grpc

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"playground/app/repository"
	"playground/app/usecase"
	"playground/config"
	"playground/driver/grpc/gen"
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
	tm, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	ur := repository.NewUserRepository(conn)
	uu := usecase.NewUserUsecase(ur, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)

	ctrl := NewController(uu)
	svr := grpc.NewServer()
	gen.RegisterPlaygroundServer(svr, ctrl)
	reflection.Register(svr)

	return svr, nil
}
