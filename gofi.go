package gofi

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
)

func isExecutable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func startedFromTerminal() bool {
	return os.Getenv("TERM") != "" && os.Getenv("FORCE_DESKTOP") != "true"
}

func empty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}

func getShell() string {
	shell := os.Getenv("SHELL")

	if len(shell) == 0 {
		shell = "sh"
	}

	return shell
}

type Executable struct {
	Name    string
	Options string
	Desktop bool
}

type GofiOptions struct {
	Executables  []Executable
	ForceDesktop bool
	Description  string
}

func (g *GofiOptions) Validate() error {
	if !startedFromTerminal() {
		g.ForceDesktop = true
	}

	if len(g.Executables) == 0 {
		if isExecutable("fzf") {
			g.Executables = append(g.Executables, Executable{
				Name:    "fzf",
				Options: "-m -i",
				Desktop: false,
			})
		}

		if isExecutable("rofi") {
			var options string

			if g.Description != "" {
				options = fmt.Sprintf("rofi -dmenu -multi-select -matching fuzzy -i -p '%s'", g.Description)
			} else {
				options = "-dmenu -multi-select -matching fuzzy -i"
			}

			g.Executables = append(g.Executables, Executable{
				Name:    "rofi",
				Options: options,
				Desktop: true,
			})
		}

		if isExecutable("dmenu") {
			g.Executables = append(g.Executables, Executable{
				Name:    "dmenu",
				Options: "",
				Desktop: true,
			})
		}
	}

	if len(g.Executables) == 0 {
		return errors.New("No executables found in PATH")
	}

	return nil
}

func (g *GofiOptions) Executable() (error, string) {
	for _, e := range g.Executables {
		if g.ForceDesktop == e.Desktop {
			return nil, e.Name + " " + e.Options
		}
	}

	return errors.New("No suitable executables found"), ""
}

func FromMap(opt *GofiOptions, input map[string]string) (error, []string) {
	rs := []string{}

	err := opt.Validate()
	if err != nil {
		return err, rs
	}

	err, command := opt.Executable()
	if err != nil {
		return err, rs
	}

	cmd := exec.Command(getShell(), "-c", command)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err, rs
	}

	go func() {
		for key := range input {
			fmt.Fprintln(stdin, key)
		}
		stdin.Close()
	}()

	stdout, err := cmd.Output()
	if err != nil {
		return err, rs
	}

	for _, key := range strings.Split(string(stdout), "\n") {
		if !empty(key) {
			rs = append(rs, strings.TrimSpace(input[key]))
		}
	}

	return nil, rs
}

func FromFilter(opt *GofiOptions, input func(in io.WriteCloser)) (error, []string) {
	rs := []string{}

	err := opt.Validate()
	if err != nil {
		return err, rs
	}

	err, command := opt.Executable()
	if err != nil {
		return err, rs
	}

	cmd := exec.Command(getShell(), "-c", command)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err, rs
	}

	go func() {
		input(stdin)
		stdin.Close()
	}()

	stdout, err := cmd.Output()
	if err != nil {
		return err, rs
	}

	for _, key := range strings.Split(string(stdout), "\n") {
		if !empty(key) {
			rs = append(rs, strings.TrimSpace(key))
		}
	}

	return nil, rs
}
