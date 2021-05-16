package database

import (
	"database/sql"
	"fmt"
	"reflect"
	"strings"
)

func intsToJson(values []int) string {
	return strings.Replace(fmt.Sprint(values), " ", ",", -1)
}

func int64sToJson(values []int64) string {
	return strings.Replace(fmt.Sprint(values), " ", ",", -1)
}

type tableModel interface {
	ptrTo(column string) interface{}
}

// returns a slice of model pointers
var runQuery = func(tx *sql.Tx, modelType reflect.Type, sql string, args ...interface{}) (interface{}, error) {
	rows, err := tx.Query(sql, args...)
	if err != nil {
		return nil, err
	}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	models := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(modelType)), 0, 0)
	for rows.Next() {
		m := reflect.New(modelType).Interface()
		values := make([]interface{}, len(columns))
		for i, column := range columns {
			values[i] = m.(tableModel).ptrTo(column)
		}
		if err = rows.Scan(values...); err != nil {
			return models, err
		}
		models = reflect.Append(models, reflect.ValueOf(m))
	}
	return models.Interface(), nil
}

func execUpdate(tx *sql.Tx, sql string, args ...interface{}) (sql.Result, error) {
	stmt, err := tx.Prepare(sql)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	return stmt.Exec(args...)
}

var runUpdate = func(tx *sql.Tx, sql string, args ...interface{}) (int64, error) {
	rs, err := execUpdate(tx, sql, args...)
	if err != nil {
		return 0, err
	}
	return rs.RowsAffected()
}

var runInsert = func(tx *sql.Tx, sql string, args ...interface{}) (int64, error) {
	rs, err := execUpdate(tx, sql, args...)
	if err != nil {
		return 0, err
	}
	return rs.LastInsertId()
}
