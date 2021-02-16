package gofi

import (
	"os"
	"testing"
)

func TestIsExecutable(t *testing.T) {
	binaries := []string{"fzf", "rofi", "dmenu"}

	for _, binary := range binaries {
		if !isExecutable(binary) {
			t.Errorf("Could not find %s on the path", binary)
		}
	}
}

func TestCustomShell(t *testing.T) {
	os.Setenv("SHELL", "bash")

	if getShell() != "bash" {
		t.Errorf("Could not find bash on the path")
	}
}

func TestDefaultShell(t *testing.T) {
	shell := os.Getenv("SHELL")

	if getShell() != shell {
		t.Errorf("Could not find %s on the path", shell)
	}
}

type bla struct {
	S string
}

func (r bla) Value() string {
	return r.S
}

func test(b Option) string {
	return b.Value()
}

func TestBla(t *testing.T) {
	p := map[string]bla{
		"a": {"f"},
		"b": {"d"},
	}

	t.Errorf("Could not find %s on the path", test(p["a"]))
}
