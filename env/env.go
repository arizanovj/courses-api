package env

import (
	"database/sql"

	"gopkg.in/doug-martin/goqu.v4"
)

type Env struct {
	DB         *sql.DB
	QB         *goqu.Database
	SQLDialect string
	BaseDir    string
	ImageDir   string
	VideoDir   string
	AppURL     string
}
