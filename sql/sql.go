package sql

import (
	"database/sql/driver"
	"time"
)

type Translated interface {
}

type Conn interface {
	Statement(translatedSql Translated) Stmt
	TranslateStatement(sql string, columns ...interface{}) Stmt
	Close() error
	Exec(translatedSql Translated, inputs ...driver.Value) (driver.Result, error)
}

type Stmt interface {
	Close() error
	Query(inputs ...driver.Value) (Rows, error)
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

type Rows interface {
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

var Translate func(sql string, columns ...interface{}) Translated
