package file

import (
	"strings"
	"testing"
)

func TestReadLine(t *testing.T) {
	err := ReadLine("./file.go", func(line string) bool {
		t.Log(line)
		if strings.Contains(line, "require") {
			return false
		}
		return true
	})
	t.Log(err)
}
