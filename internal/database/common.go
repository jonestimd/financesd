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
	PtrTo(column string) interface{}
}

// returns a slice of model pointers
var runQuery = func(tx *sql.Tx, modelType reflect.Type, sql string, args ...interface{}) interface{} {
	rows, err := tx.Query(sql, args...)
	if err != nil {
		panic(err)
	}
	columns, err := rows.Columns()
	if err != nil {
		panic(err)
	}
	models := reflect.MakeSlice(reflect.SliceOf(reflect.PtrTo(modelType)), 0, 0)
	for rows.Next() {
		m := reflect.New(modelType).Interface()
		values := make([]interface{}, len(columns))
		for i, column := range columns {
			values[i] = m.(tableModel).PtrTo(column)
		}
		if err = rows.Scan(values...); err != nil {
			panic(err)
		}
		models = reflect.Append(models, reflect.ValueOf(m))
	}
	return models.Interface()
}

func execUpdate(tx *sql.Tx, sql string, args ...interface{}) sql.Result {
	stmt, err := tx.Prepare(sql)
	if err != nil {
		panic(err)
	}
	defer stmt.Close()
	result, err := stmt.Exec(args...)
	if err != nil {
		panic(err)
	}
	return result
}

var runUpdate = func(tx *sql.Tx, sql string, args ...interface{}) int64 {
	count, err := execUpdate(tx, sql, args...).RowsAffected()
	if err != nil {
		panic(err)
	}
	return count
}

var runInsert = func(tx *sql.Tx, sql string, args ...interface{}) int64 {
	id, err := execUpdate(tx, sql, args...).LastInsertId()
	if err != nil {
		panic(err)
	}
	return id
}
