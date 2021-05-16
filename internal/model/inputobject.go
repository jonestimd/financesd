package model

import "errors"

type inputObject map[string]interface{}

func (io inputObject) getInt(key string) (interface{}, bool) {
	if value, exists := io[key]; exists {
		if i, ok := value.(int); ok {
			return int64(i), true
		}
		return nil, true
	}
	return nil, false
}

func (io inputObject) intOrNull(key string) interface{} {
	if value, ok := io[key].(int); ok {
		return int64(value)
	}
	return nil
}

func (io inputObject) requireInt(key string) int64 {
	return int64(io[key].(int))
}

func (io inputObject) getFloat(key string) (interface{}, bool) {
	if value, exists := io[key]; exists {
		if f, ok := value.(float64); ok {
			return f, true
		}
		return nil, true
	}
	return nil, false
}

func (io inputObject) floatOrNull(key string) interface{} {
	if value, ok := io[key].(float64); ok {
		return value
	}
	return nil
}

func (io inputObject) getString(key string) (interface{}, bool) {
	if value, exists := io[key]; exists {
		if f, ok := value.(string); ok {
			return f, true
		}
		return nil, true
	}
	return nil, false
}

func (io inputObject) stringOrNull(key string) interface{} {
	if s, ok := io[key].(string); ok {
		return s
	}
	return nil
}

func (io inputObject) getVersionID() (*versionID, error) {
	if id, ok := io.getInt("id"); ok {
		if version, ok := io.getInt("version"); ok {
			return &versionID{ID: id.(int64), Version: version.(int64)}, nil
		}
		return nil, errors.New("version is required for update/delete")
	}
	return nil, nil
}
