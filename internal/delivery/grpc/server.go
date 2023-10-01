package grpc

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"playground/internal/config"
	"playground/internal/delivery/grpc/gen"
	"playground/internal/pkg/token"
	"playground/internal/wallet/mq"
	"playground/internal/wallet/repository"
	"playground/internal/wallet/usecase"
)

// NewServer creates a new gRPC server.
func NewServer(cfg config.Config) (*grpc.Server, error) {
	conn, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%d)/playground?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort))
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
	p := mq.NewAsynqProducer(asynq.RedisClientOpt{
		Addr: cfg.RedisAddress,
	})
	r := repository.NewRepository(conn)
	u := usecase.NewUsecase(r, p, nil, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	ctrl := NewController(u, tm)
	svr := grpc.NewServer(grpc.UnaryInterceptor(GRPCLogger))
	gen.RegisterPlaygroundServer(svr, ctrl)
	reflection.Register(svr)

	return svr, nil
}
