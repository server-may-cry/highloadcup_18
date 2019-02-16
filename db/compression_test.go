package db

import "testing"

func TestCompressionStorage(t *testing.T) {
	compressor := newCompressionStorage()
	test := compressor.getIntOrCreate("test")
	test2 := compressor.getIntOrCreate("test2")
	test3 := compressor.getIntOrCreate("test")
	if test != test3 {
		t.Errorf("Expected same ID %d %d", test, test3)
	}
	if test == test2 {
		t.Errorf("Expected different ID %d %d", test, test2)
	}
	s := compressor.getString(test)
	if s != "test" {
		t.Errorf("Expected string 'test' got '%s'", s)
	}
	s2 := compressor.getString(test2)
	if s2 != "test2" {
		t.Errorf("Expected string 'test2' got '%s'", s)
	}
}

func TestStatus(t *testing.T) {
	s := StatusFreeInt
	str := intToStatus[s]
	r := statusToInt[str]
	if s != r {
		t.Errorf("Expected same %d %d", s, r)
	}
}
