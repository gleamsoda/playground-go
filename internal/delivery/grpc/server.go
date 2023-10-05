package grpc

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net"
	"os/signal"
	"syscall"

	_ "github.com/go-sql-driver/mysql"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"playground/internal/config"
	"playground/internal/delivery/grpc/gen"
	"playground/internal/delivery/grpc/handler"
	"playground/internal/delivery/grpc/interceptor"
	"playground/internal/pkg/token"
	"playground/internal/wallet/dispatcher"
	"playground/internal/wallet/repository"
	"playground/internal/wallet/usecase"
)

type Server struct {
	listener net.Listener
	server   *grpc.Server
}

func Run(ctx context.Context) error {
	cfg := config.Get()

	// run gRPC server
	server, err := NewServer(cfg)
	if err != nil {
		return err
	}
	ctxn, canceln := signal.NotifyContext(ctx, syscall.SIGINT, syscall.SIGTERM)
	defer canceln()
	ctxs, cancels := context.WithCancelCause(ctx)
	defer cancels(nil)
	go func() {
		log.Info().Msgf("start gRPC server at %s", server.listener.Addr())
		if err := server.Run(); err != nil {
			cancels(err)
		}
		cancels(nil)
	}()

	// run gateway server
	gw, err := NewGatewayServer(ctxs, cfg)
	if err != nil {
		return err
	}
	ctxg, cancelg := context.WithCancelCause(ctxn)
	defer cancelg(nil)
	go func() {
		log.Info().Msgf("start gateway server at %s", gw.server.Addr)
		if err := gw.Run(); err != nil {
			cancelg(err)
		}
		cancelg(nil)
	}()

	select {
	case <-ctxn.Done():
		log.Info().Msg("shutting down gateway server gracefully...")
		if err := gw.Shutdown(); err != nil {
			log.Error().Err(err)
		} else {
			log.Info().Msg("shutdown gateway server complete")
		}
		log.Info().Msg("shutting down gRPC server gracefully...")
		server.Shutdown()
		log.Info().Msg("shutdown gRPC server complete")
	case <-ctxs.Done():
		log.Info().Msg("shutting down gateway server gracefully...")
		if err := gw.Shutdown(); err != nil {
			log.Error().Err(err)
		} else {
			log.Info().Msg("shutdown gateway server complete")
		}
		if cause := context.Cause(ctxs); !errors.Is(cause, context.Canceled) {
			return fmt.Errorf("gRPC server error: %v", cause)
		}
	case <-ctxg.Done():
		log.Info().Msg("shutting down gRPC server gracefully...")
		server.Shutdown()
		log.Info().Msg("shutdown gRPC server complete")
		if cause := context.Cause(ctxg); !errors.Is(cause, context.Canceled) {
			return fmt.Errorf("gateway server error: %v", cause)
		}
	}
	return nil
}

// NewServer creates a new gRPC server.
func NewServer(cfg config.Config) (*Server, error) {
	listener, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		return nil, err
	}
	tm, err := token.NewPasetoManager(cfg.TokenSymmetricKey)
	if err != nil {
		return nil, fmt.Errorf("cannot create token maker: %w", err)
	}
	// repositories
	conn, err := sql.Open("mysql", cfg.DBName())
	if err != nil {
		return nil, err
	} else if err := conn.Ping(); err != nil {
		return nil, err
	}
	r := repository.NewRepository(conn)
	p := dispatcher.NewDispatcher(asynq.RedisClientOpt{
		Addr: cfg.RedisAddress,
	})
	// usecases
	u := usecase.NewUsecase(r, p, nil, tm, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)

	ctrl := handler.NewHandler(u, tm)
	server := grpc.NewServer(grpc.ChainUnaryInterceptor(
		interceptor.Logger,
		interceptor.ErrorHandler,
	))
	gen.RegisterPlaygroundServer(server, ctrl)
	reflection.Register(server)

	return &Server{
		server:   server,
		listener: listener,
	}, nil
}

func (s *Server) Run() error {
	return s.server.Serve(s.listener)
}

func (s *Server) Shutdown() {
	s.server.GracefulStop()
}
