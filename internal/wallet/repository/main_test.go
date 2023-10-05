package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"playground/internal/config"
	"playground/internal/wallet"
)

var repository wallet.Repository

func TestMain(m *testing.M) {
	cfg := config.Get()
	conn, err := sql.Open("mysql", cfg.DBName())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else if err := conn.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	repository = NewRepository(conn)
	os.Exit(m.Run())
}
