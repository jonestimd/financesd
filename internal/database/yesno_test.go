package database

import (
	"reflect"
	"testing"
)

func Test_YesNo_Value(t *testing.T) {
	yn := YesNo('Y')

	value, _ := yn.Value()

	if !reflect.DeepEqual(value, []byte{'Y'}) {
		t.Errorf("Expected: %v, got: %v", yn, value)
	}
}

func Test_YesNo_Scan(t *testing.T) {
	yn := YesNo(0)

	yn.Scan([]byte{'Y'})

	if yn != 'Y' {
		t.Errorf("Expected: 1, got: %v", yn)
	}
}

func Test_YesNo_Get(t *testing.T) {
	y := YesNo('Y')
	if !y.Get() {
		t.Error("Expected true")
	}

	n := YesNo('N')
	if n.Get() {
		t.Error("Expected false")
	}
}

func Test_YesNo_Set(t *testing.T) {
	yn := YesNo(0)

	yn.Set(false)
	if byte(yn) != 'N' {
		t.Errorf("Expected: N, got: %v", yn)
	}

	yn.Set(true)
	if byte(yn) != 'Y' {
		t.Errorf("Expected: Y, got: %v", yn)
	}
}
