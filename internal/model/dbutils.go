package model

import (
	"database/sql"
	"reflect"
)

type model interface {
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
			values[i] = m.(model).ptrTo(column)
		}
		if err = rows.Scan(values...); err != nil {
			return models, err
		}
		models = reflect.Append(models, reflect.ValueOf(m))
	}
	return models.Interface(), nil
}

type idSet struct {
	ids map[int64]struct{}
}

func newIDSet() *idSet {
	return &idSet{ids: make(map[int64]struct{})}
}

func (s *idSet) Add(id int64) {
	if _, ok := s.ids[id]; !ok {
		s.ids[id] = struct{}{}
	}
}

func (s *idSet) Values() []int64 {
	ids := make([]int64, len(s.ids))
	i := 0
	for id := range s.ids {
		ids[i] = id
		i++
	}
	return ids
}
