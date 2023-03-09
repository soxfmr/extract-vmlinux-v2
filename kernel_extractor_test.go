package vmlinux

import (
	"fmt"
	"os"
	"testing"
)

func TestExtractTo(t *testing.T) {
	file, err := os.Open("/tmp/vmlinuz-5.10.0-19-amd64")
	if err != nil {
		panic(err)
	}

	tmp, err := os.CreateTemp("", "vmlinux")
	if err != nil {
		panic(err)
	}

	if err := ExtractTo(file, tmp); err != nil {
		panic(err)
	}

	// os.RemoveAll(tmp.Name())

	fmt.Println(tmp.Name())
}
