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

var rm app.RepositoryManager

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
	do.Provide(injector, NewManager)
	do.ProvideValue(injector, conn)
	rm = do.MustInvoke[app.RepositoryManager](injector)

	os.Exit(m.Run())
}
