package utils

import (
	"io"
	"testing"
)

func TestMessageWriter(t *testing.T) {
	out := NewMessageWriter()

	exp := "Test"

	writerMock(out, exp)
	t.Log(out)

	act := out.Read()
	if exp != act {
		t.Fatalf("unexpected read data '%s' expect '%s'", act, exp)
	}
}

func writerMock(out io.Writer, data string) {
	out.Write([]byte(data))
}
