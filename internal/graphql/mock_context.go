package graphql

import (
	"database/sql"
	"time"
)

type mockSql struct {
	method string
	query string
	args []interface{}
}

func (sql *mockSql) Exec(query string, args ...interface{}) (sql.Result, error) {
	sql.method = "Exec"
	sql.query = query
	sql.args = args
	return nil, nil
}

func (sql *mockSql) Prepare(query string) (*sql.Stmt, error) {
	sql.method = "Prepare"
	sql.query = query
return nil, nil
}

func (sql *mockSql) Query(query string, args ...interface{}) (*sql.Rows, error) {
	sql.method = "Query"
	sql.query = query
	sql.args = args
	return nil, nil
}

func (sql *mockSql) QueryRow(query string, args ...interface{}) *sql.Row {
	sql.method = "QueryRow"
	sql.query = query
	sql.args = args
	return nil
}

type mockContext struct {
	sqlCommon *mockSql
}

func newMockContext() *mockContext {
	return &mockContext{sqlCommon: &mockSql{}}
}

func (ctx *mockContext) Value(key interface{}) interface{} {
	if key == DbContextKey {
		return ctx.sqlCommon
	}
	return nil
}

func (ctx *mockContext) Deadline() (deadline time.Time, ok bool) {
	panic("unexpected call to Deadline")
}

func (ctx *mockContext) Done() <-chan struct{} {
	panic("unexpected call to Done")
}

func (ctx *mockContext) Err() error {
	panic("unexpected call to Err")
}
