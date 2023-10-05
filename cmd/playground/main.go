package main

import (
	"fmt"
	"os"

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
	cmd := NewCmdRoot()
	cmd.AddCommand(NewCmdAsynq(), NewCmdGin(), NewCmdGRPC(), NewCmdMigrateUp())
	if err := cmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func NewCmdRoot() *cobra.Command {
	cmd := &cobra.Command{
		Use:           "playground",
		SilenceUsage:  true,
		SilenceErrors: true,
		PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
			if config.Get().IsDevelopment() {
				log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
			}
			return nil
		},
	}
	return cmd
}

func NewCmdGin() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gin",
		Short: "Run gin server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return gin.Run(cmd.Context())
		},
	}
	return cmd
}

func NewCmdGRPC() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "grpc",
		Short: "Run gRPC server",
		RunE: func(cmd *cobra.Command, args []string) error {
			return grpc.Run(cmd.Context())
		},
	}
	return cmd
}

func NewCmdAsynq() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "asynq",
		Short: "Run asynq",
		RunE: func(cmd *cobra.Command, args []string) error {
			return asynq.Run()
		},
	}
	return cmd
}

func NewCmdMigrateUp() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "migrate-up",
		Short: "Run database migration up",
		RunE: func(_ *cobra.Command, _ []string) error {
			cfg := config.Get()
			migration, err := migrate.New(cfg.MigrationURL, fmt.Sprintf("mysql://%s", cfg.DBName()))
			if err != nil {
				log.Fatal().Err(err).Msg("cannot create new migrate instance")
			}
			if err = migration.Up(); err != nil && err != migrate.ErrNoChange {
				log.Fatal().Err(err).Msg("failed to run migrate up")
			}
			log.Info().Msg("db migrated successfully")
			return nil
		},
	}
	return cmd
}
