package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"

	"playground/config"
	"playground/driver/gin"
	"playground/driver/grpc"
)

func main() {
	if err := NewCmdRoot().Execute(); err != nil {
		log.Fatal(err)
	}
}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "playground",
		SilenceUsage:  true,
		SilenceErrors: true,
	}
	cmd.AddCommand(NewCmdGin(), NewCmdGRPC())
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

func runGin(cmd *cobra.Command, args []string) error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	srv, err := gin.NewServer(cfg)
	if err != nil {
		return err
	}
	ctx, stop := signal.NotifyContext(cmd.Context(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	errCh := make(chan error, 1)
	go func() {
		defer close(errCh)
		if err := srv.ListenAndServe(); errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		log.Println("shutting down gracefully...")
		shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Minute)
		defer shutdownCancel()
		if err := srv.Shutdown(shutdownCtx); err != nil {
			return fmt.Errorf("failed to shutdown gracefully: %v", err)
		}
		log.Println("shutdown complete")
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
		log.Printf("start gRPC server at %s", l.Addr())
		if err := svr.Serve(l); err != nil {
			serviceErrCh <- err
		}
	}()
	gatewayErrCh := make(chan error, 1)
	go func() { // run gateway server
		defer close(gatewayErrCh)
		log.Printf("start gateway server at %s", gw.Addr)
		if err := gw.ListenAndServe(); err != nil {
			gatewayErrCh <- err
		}
	}()

	select {
	case <-ctx.Done():
		if err := shutdownGateway(gw); err != nil {
			log.Println(err)
		}
		log.Println("shutting down gRPC server gracefully...")
		svr.GracefulStop()
		log.Println("shutdown gRPC server complete")
	case err := <-serviceErrCh:
		if err := shutdownGateway(gw); err != nil {
			log.Println(err)
		}
		return fmt.Errorf("gRPC server error: %v", err)
	case gatewayErr := <-gatewayErrCh:
		log.Println("shutting down gRPC server gracefully...")
		svr.GracefulStop()
		log.Println("shutdown gRPC server complete")
		return fmt.Errorf("gateway server error: %v", gatewayErr)
	}
	return nil
}

func shutdownGateway(gw *http.Server) error {
	log.Println("shutting down gateway server gracefully...")
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer shutdownCancel()
	if err := gw.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf("gateway server error: failed to shutdown gracefully: %v", err)
	}
	log.Println("shutdown gateway server complete")
	return nil
}
