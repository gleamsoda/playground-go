package main

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"playground/internal/config"
	"playground/internal/delivery/asynq"
	"playground/internal/delivery/gin"
	"playground/internal/delivery/grpc"
)

func main() {
	if err := NewCmdRoot().Execute(); err != nil {
		log.Fatal().Err(err)
	}
}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "playground",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(NewCmdGin(), NewCmdGRPC(), NewCmdAsynq())
	return cmd
}

func NewCmdGin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gin",
		Short: "Run gin server",
		RunE:  runGin,
	}
	return cmd
}

func NewCmdGRPC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grpc",
		Short: "Run gRPC server",
		RunE:  runGRPC,
	}
	return cmd
}

func NewCmdAsynq() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asynq",
		Short: "Run asynq",
		RunE:  runAsynq,
	}
	return cmd
}

func runAsynq(cmd *cobra.Command, args []string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	runDBMigration(cfg)

	srv, err := asynq.NewConsumer(cfg)
	if err != nil {
		return err
	}

	log.Info().Msg("start task processor")
	srv.Start()
	return nil
}

func runGin(cmd *cobra.Command, args []string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	runDBMigration(cfg)
	srv, err := gin.NewServer(cfg)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		if err := srv.Start(); errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Info().Msg("shutting down gracefully...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer shutdownCancel()
		if err := srv.Stop(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown gracefully: %v", err)
		}
		log.Info().Msg("shutdown complete")
	case err := <-errCh:
		return fmt.Errorf("error: %v", err)
	}
	return nil
}

func runGRPC(cmd *cobra.Command, args []string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	if cfg.Environment == "development" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}
	runDBMigration(cfg)
	svr, err := grpc.NewServer(cfg)
	if err != nil {
		return err
	}
	l, err := net.Listen("tcp", cfg.GRPCServerAddress)
	if err != nil {
		return err
	}
	gw, err := grpc.NewGatewayServer(cmd.Context(), cfg)
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	serviceErrCh := make(chan error, 1)
	go func() { // run gRPC server
		defer close(serviceErrCh)
		log.Info().Msgf("start gRPC server at %s", l.Addr())
		if err := svr.Serve(l); err != nil {
			serviceErrCh <- err
		}
	}()
	gatewayErrCh := make(chan error, 1)
	go func() { // run gateway server
		defer close(gatewayErrCh)
		log.Info().Msgf("start gateway server at %s", gw.Addr)
		if err := gw.ListenAndServe(); err != nil {
			gatewayErrCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		if err := shutdownGateway(gw); err != nil {
			log.Error().Err(err)
		}
		log.Info().Msg("shutting down gRPC server gracefully...")
		svr.GracefulStop()
		log.Info().Msg("shutdown gRPC server complete")
	case err := <-serviceErrCh:
		if err := shutdownGateway(gw); err != nil {
			log.Error().Err(err)
		}
		return fmt.Errorf("gRPC server error: %v", err)
	case gatewayErr := <-gatewayErrCh:
		log.Info().Msg("shutting down gRPC server gracefully...")
		svr.GracefulStop()
		log.Info().Msg("shutdown gRPC server complete")
		return fmt.Errorf("gateway server error: %v", gatewayErr)
	}
	return nil
}

func shutdownGateway(gw *http.Server) error {
	log.Info().Msg("shutting down gateway server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer shutdownCancel()
	if err := gw.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("gateway server error: failed to shutdown gracefully: %v", err)
	}
	log.Info().Msg("shutdown gateway server complete")
	return nil
}

func runDBMigration(cfg config.Config) {
	migration, err := migrate.New(cfg.MigrationURL, fmt.Sprintf("mysql://%s:%s@tcp(%s:%d)/playground?parseTime=true", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort))
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create new migrate instance")
	}

	if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
		log.Fatal().Err(err).Msg("failed to run migrate up")
	}

	log.Info().Msg("db migrated successfully")
}
