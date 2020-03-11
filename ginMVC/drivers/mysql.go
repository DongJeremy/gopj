package drivers

import (
	"fmt"
	"time"

	"github.com/davidddw/gopj/ginMVC/configs"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
)

// query need rows.Close to release db ins
// exec will release automatic

func InitMySQL() (MySQLDB *xorm.Engine, MySQLDBErr error) {
	// get db config
	dbConfig := configs.SystemConfig.MySQLConfig

	dbDSN := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s",
		dbConfig.User,
		dbConfig.Passwd,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Dbname,
		dbConfig.Charset,
	)

	// connect and open db connection
	MySQLDB, MySQLDBErr = xorm.NewEngine("mysql", dbDSN)
	MySQLDB.ShowSQL(true)

	if MySQLDBErr != nil {
		panic("database data source name error: " + MySQLDBErr.Error())
	}

	// max open connections
	MySQLDB.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)

	// max idle connections
	MySQLDB.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)

	// max lifetime of connection if <=0 will forever
	MySQLDB.DB().SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetimeConns))

	// check db connection at once avoid connect failed
	// else error will be reported until db first sql operate
	if MySQLDBErr = MySQLDB.DB().Ping(); nil != MySQLDBErr {
		panic("database connect failed: " + MySQLDBErr.Error())
	}
	return
}
