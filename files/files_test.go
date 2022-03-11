package files

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
)

func Test_read(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Fatalf(err.Error())
	}
	fp := filepath.Join(home, ".aria2", "aria2.conf")
	data, err := ReadFile(fp)
	if err != nil {
		t.Fatalf(err.Error())
	}
	fmt.Println(string(data))
}
