package database

import "errors"

type InputObject map[string]interface{}

func (io InputObject) GetInt(key string) (interface{}, bool) {
	if value, exists := io[key]; exists {
		if i, ok := value.(int); ok {
			return int64(i), true
		}
		return nil, true
	}
	return nil, false
}

func (io InputObject) IntOrNull(key string) interface{} {
	if value, ok := io[key].(int); ok {
		return int64(value)
	}
	return nil
}

func (io InputObject) RequireInt(key string) int64 {
	return int64(io[key].(int))
}

func (io InputObject) GetFloat(key string) (interface{}, bool) {
	if value, exists := io[key]; exists {
		if f, ok := value.(float64); ok {
			return f, true
		}
		return nil, true
	}
	return nil, false
}

func (io InputObject) FloatOrNull(key string) interface{} {
	if value, ok := io[key].(float64); ok {
		return value
	}
	return nil
}

func (io InputObject) GetString(key string) (interface{}, bool) {
	if value, exists := io[key]; exists {
		if f, ok := value.(string); ok {
			return f, true
		}
		return nil, true
	}
	return nil, false
}

func (io InputObject) StringOrNull(key string) interface{} {
	if s, ok := io[key].(string); ok {
		return s
	}
	return nil
}

type VersionID struct {
	ID      int64
	Version int64
}

func (io InputObject) GetVersionID() (*VersionID, error) {
	if id, ok := io.GetInt("id"); ok {
		if version, ok := io.GetInt("version"); ok {
			return &VersionID{ID: id.(int64), Version: version.(int64)}, nil
		}
		return nil, errors.New("version is required for update/delete")
	}
	return nil, nil
}
