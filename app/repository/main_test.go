package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"

	"playground/app"
	"playground/config"
)

var repository app.Repository

func TestMain(m *testing.M) {
	cfg, err := config.NewConfig()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	conn, err := sql.Open("mysql", cfg.DBSource)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	if err := conn.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	repository = NewRepository(conn)
	os.Exit(m.Run())
}
