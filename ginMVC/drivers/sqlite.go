package drivers

import (
	"github.com/davidddw/gopj/ginMVC/configs"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-sqlite3"
)

func InitSqlite() (SqliteDB *xorm.Engine, SqliteDBErr error) {
	// get db config
	dbConfig := configs.SystemConfig.SqliteConfig

	// connect and open db connection
	SqliteDB, SqliteDBErr = xorm.NewEngine("sqlite3", dbConfig.Filename)
	SqliteDB.ShowSQL(true)

	if SqliteDBErr != nil {
		panic("database data source name error: " + SqliteDBErr.Error())
	}

	// check db connection at once avoid connect failed
	// else error will be reported until db first sql operate
	if SqliteDBErr = SqliteDB.DB().Ping(); nil != SqliteDBErr {
		panic("database connect failed: " + SqliteDBErr.Error())
	}
	return
}
