package main

import (
	"os"
	"testing"
)

func TestDebugSimple(t *testing.T) {
	expextedOutput := "value: 0"
	output := captureOutput(1000, func() { debug("value:", 0) })
	if expextedOutput != output {
		t.Errorf("Expected [%s] but was [%s]", expextedOutput, output)
	}
}

func TestDebugDataSet(t *testing.T) {
	tests := []struct {
		expected   string
		bufferSize int
		values     []any
	}{
		{"value: 0", 20, []any{"value:", 0}},
		{"{ value1: 0 , value2: toto }", 30, []any{"{", "value1:", 0, ", value2:", "toto", "}"}},
	}
	for _, test := range tests {
		testName := "Debug " + test.expected
		t.Run(testName, func(t *testing.T) {
			output := captureOutput(test.bufferSize, func() { debug(test.values...) })
			if test.expected != output {
				t.Errorf("Expected [%s] but was [%s]", test.expected, output)
			}
		})
	}
}

func BenchmarkDebug(b *testing.B) {
	tests := []struct {
		expected   string
		bufferSize int
		values     []any
	}{
		{"value: 0", 20, []any{"value:", 0}},
		{"{ value1: 0 , value2: toto }", 30, []any{"{", "value1:", 0, ", value2:", "toto", "}"}},
	}
	nbTests := len(tests)
	for i := 0; i < b.N; i++ {
		test := tests[i%nbTests]
		captureOutput(test.bufferSize, func() {
			debug(test.values...)
		})
	}
}

func captureOutput(bufferSize int, function func()) string {
	stderr := os.Stderr
	defer func() { os.Stderr = stderr }()
	read, write, _ := os.Pipe()
	defer write.Close()
	defer read.Close()
	os.Stderr = write
	function()
	bytes := make([]byte, bufferSize)
	n, err := read.Read(bytes)
	if err != nil {
		panic("Capture output error")
	}
	if n == bufferSize {
		panic("Capture output buffer full")
	}
	return string(bytes[:n-1]) // remove new line
}
