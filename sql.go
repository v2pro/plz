package plz

import (
	"database/sql/driver"
	"time"
)

var SqlOpen func(drv driver.Driver, dsn string) (SqlConn, error)

type TranslatedSql interface {
}

type SqlConn interface {
	Statement(translatedSql TranslatedSql) SqlStmt
	TranslateStatement(sql string, columns ...interface{}) SqlStmt
	Close() error
}

type SqlStmt interface {
	Close() error
	Query(inputs ...driver.Value) (SqlRows, error)
	Exec(inputs ...driver.Value) (driver.Result, error)
}

type ColumnIndex int

type ColumnGroup struct {
	Group                string
	Columns              []string
	BatchInsertRowsCount int
}

func Columns(group string, columns ...string) ColumnGroup {
	return ColumnGroup{group, columns, 0}
}

func BatchInsertColumns(batchInsertRowsCount int, columns ...string) ColumnGroup {
	return ColumnGroup{"COLUMNS", columns, batchInsertRowsCount}
}

type SqlRows interface {
	Close() error
	Next() error
	C(column string) ColumnIndex
	GetString(idx ColumnIndex) string
	GetInt64(idx ColumnIndex) int64
	GetInt(idx ColumnIndex) int
	GetTime(idx ColumnIndex) time.Time
	GetByteArray(idx ColumnIndex) []byte
}

func BatchInsertRow(inputs ...driver.Value) []driver.Value {
	return inputs
}

var TranslateSql func(sql string, columns ...interface{}) TranslatedSql
