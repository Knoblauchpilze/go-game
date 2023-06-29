package example

import "testing"

// test function
func TestPrintHelloWorld(t *testing.T) {
	actualString := FormatHelloWorld("toto")
	expectedString := "Hello world from toto!\n"
	if actualString != expectedString {
		t.Errorf("Expected String(%s) is not same as"+
			" actual string (%s)", expectedString, actualString)
	}
}
