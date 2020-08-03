package conf

import (
	"database/sql"

	"github.com/jinzhu/gorm"

	"tapera/util/env"
)

// GormDb func
func (f *Factory) GormDb(cfg *DbConfig, initf func()) *gorm.DB {
	db := f.Db(cfg)
	return gormDb(db, cfg.Dialect, initf)
}

// GormDb func
func (f *EnvFactory) GormDb(initf func()) *gorm.DB {
	db := f.Db()
	dialect := env.Str(envDbDialect, false, "postgres")
	return gormDb(db, dialect, initf)
}

func gormDb(db *sql.DB, dialect string, initf func()) *gorm.DB {
	gdb, err := gorm.Open(dialect, db)
	if err != nil {
		panic(err)
	}

	if initf != nil {
		initf()
	}

	return gdb
}
