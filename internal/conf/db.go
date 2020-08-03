package conf

import (
	"database/sql"

	"tapera/util/env"
)

const (
	envDbConnStr       = "DB_CONN_STR"
	envDbConnStrAsFile = "DB_CONN_STR_AS_FILE"
	envDbMaxIdleConn   = "DB_MAX_IDLE_CONN"
	envDbMaxOpenConn   = "DB_MAX_OPEN_CONN"
	envDbDialect       = "DB_DIALECT"
)

// DbConfig func
type DbConfig struct {
	ConnStr       string
	ConnStrAsFile bool
	MaxIdleConn   int
	MaxOpenConn   int
	Dialect       string
}

// Db func
func (f *Factory) Db(cfg *DbConfig) *sql.DB {
	connStr := cfg.ConnStr
	if cfg.ConnStrAsFile {
		str, err := readTxtFromFile(connStr)
		if err != nil {
			panic(err)
		}
		connStr = str
	}

	db, err := sql.Open(cfg.Dialect, connStr)
	if err != nil {
		panic(err)
	}

	err = db.Ping()
	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(cfg.MaxIdleConn)
	db.SetMaxOpenConns(cfg.MaxOpenConn)
	return db
}

// Db func
func (f *EnvFactory) Db() *sql.DB {
	return f.factory.Db(&DbConfig{
		ConnStr:       env.Str(envDbConnStr, true, nil),
		ConnStrAsFile: env.Bool(envDbConnStrAsFile, false, nil),
		MaxIdleConn:   env.Int(envDbMaxIdleConn, false, 0),
		MaxOpenConn:   env.Int(envDbMaxOpenConn, false, 0),
		Dialect:       env.Str(envDbDialect, false, "postgres"),
	})
}
