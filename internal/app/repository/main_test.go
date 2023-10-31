package repository

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/samber/do"

	"playground/internal/app"
	"playground/internal/config"
)

var repository app.Repository

func TestMain(m *testing.M) {
	conn, err := sql.Open("mysql", config.Get().DBName())
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	} else if err := conn.Ping(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	injector := do.New()
	do.Provide(injector, NewRepository)
	do.ProvideValue(injector, conn)
	repository = do.MustInvoke[app.Repository](injector)

	os.Exit(m.Run())
}
