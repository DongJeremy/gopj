package drivers

import (
	"fmt"
	"time"

	"github.com/davidddw/gopj/ginMVC/configs"
	"github.com/go-xorm/xorm"
	_ "github.com/mattn/go-oci8"
)

func InitOracle() (OracleDB *xorm.Engine, OracleDBErr error) {
	// get db config
	dbConfig := configs.SystemConfig.OracleConfig

	dbDSN := fmt.Sprintf("%s/%s@%s:%d/%s",
		dbConfig.Username,
		dbConfig.Password,
		dbConfig.Host,
		dbConfig.Port,
		dbConfig.Sid)

	// connect and open db connection
	OracleDB, OracleDBErr = xorm.NewEngine("oci8", dbDSN)
	OracleDB.ShowSQL(true)

	if OracleDBErr != nil {
		panic("database data source name error: " + OracleDBErr.Error())
	}

	// max open connections
	OracleDB.DB().SetMaxOpenConns(dbConfig.MaxOpenConns)

	// max idle connections
	OracleDB.DB().SetMaxIdleConns(dbConfig.MaxIdleConns)

	// max lifetime of connection if <=0 will forever
	OracleDB.DB().SetConnMaxLifetime(time.Duration(dbConfig.MaxLifetimeConns))

	// check db connection at once avoid connect failed
	// else error will be reported until db first sql operate
	if OracleDBErr = OracleDB.DB().Ping(); nil != OracleDBErr {
		panic("database connect failed: " + OracleDBErr.Error())
	}
	return
}
