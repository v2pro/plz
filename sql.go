package plz

import (
	"database/sql/driver"
	"github.com/v2pro/plz/sql"
)

var OpenSqlConn func(drv driver.Driver, dsn string) (sql.Conn, error)
