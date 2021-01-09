package model

import "testing"

func Test_YesNo_Value(t *testing.T) {
	yn := YesNo{value: 'Y'}

	value, _ := yn.Value();

	if (value != yn.value) {
		t.Errorf("Expected: %v, got: %v", yn.value, value)
	}
}

func Test_YesNo_Scan(t *testing.T) {
	yn := YesNo{}

	yn.Scan([]byte{'Y'})

	if (yn.value != 'Y') {
		t.Errorf("Expected: 1, got: %v", yn.value)
	}
}

func Test_YesNo_Get(t *testing.T) {
	if !(&YesNo{value: 'Y'}).Get() {
		t.Error("Expected true")
	}

	if (&YesNo{value: 'N'}).Get() {
		t.Error("Expected false")
	}
}

func Test_YesNo_Set(t *testing.T) {
	yn := YesNo{}

	yn.Set(false)
	if yn.value != 'N' {
		t.Errorf("Expected: N, got: %v", yn.value)
	}

	yn.Set(true)
	if yn.value != 'Y' {
		t.Errorf("Expected: Y, got: %v", yn.value)
	}
}