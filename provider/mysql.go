package provider

import (
	"database/sql"
	"excel-receiver/config"
	"fmt"
	"net/url"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

func NewMysql() (*sqlx.DB, error) {
	cfg := config.Configuration.Mysql

	url := (&url.URL{
		User:     url.UserPassword(cfg.Username, cfg.Password),
		Host:     fmt.Sprintf("tcp(%s:%d)", cfg.Host, cfg.Port),
		Path:     cfg.Database,
		RawQuery: strings.Join(cfg.Options, "&"),
	}).String()

	url = strings.TrimLeft(url, "/")

	dbConn, err := sql.Open("mysql", url)
	if err != nil {
		return nil, err
	}
	db := sqlx.NewDb(dbConn, "mysql")
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}
