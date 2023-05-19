package gofi

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
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
	PreviewPath  string
}

func (g *GofiOptions) Validate() error {
	if !startedFromTerminal() {
		g.ForceDesktop = true
	}

	if len(g.Executables) == 0 {
		if isExecutable("/usr/local/bin/fzf") {
			var options string

			if g.PreviewPath != "" {
				options = fmt.Sprintf("-m -i --bind 'esc:become(exit)' --preview '%s'", g.PreviewPath)
			} else {
				options = "-m -i --bind 'esc:become(exit)'"
			}

			g.Executables = append(g.Executables, Executable{
				Name:    "/usr/local/bin/fzf",
				Options: options,
				Desktop: false,
			})
		}

		if isExecutable("fzf") {
			var options string

			if g.PreviewPath != "" {
				options = fmt.Sprintf("-m -i --bind 'esc:become(exit)' --preview '%s'", g.PreviewPath)
			} else {
				options = "-m -i --bind 'esc:become(exit)'"
			}

			g.Executables = append(g.Executables, Executable{
				Name:    "fzf",
				Options: options,
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

func (g *GofiOptions) Executable() (string, error) {
	for _, e := range g.Executables {
		if g.ForceDesktop == e.Desktop {
			return e.Name + " " + e.Options, nil
		}
	}

	return "", errors.New("No suitable executables found")
}

func FromMap(opt *GofiOptions, input map[string]string) ([]string, error) {
	rs := []string{}

	err := opt.Validate()
	if err != nil {
		return rs, err
	}

	command, err := opt.Executable()
	if err != nil {
		return rs, err
	}

	cmd := exec.Command(getShell(), "-c", command)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return rs, err
	}

	go func() {
		for key := range input {
			fmt.Fprintln(stdin, key)
		}
		stdin.Close()
	}()

	stdout, err := cmd.Output()
	if err != nil {
		return rs, err
	}

	for _, key := range strings.Split(string(stdout), "\n") {
		if !empty(key) {
			rs = append(rs, strings.TrimSpace(input[key]))
		}
	}

	return rs, nil
}

func FromFilter(opt *GofiOptions, input func(in io.WriteCloser)) ([]string, error) {
	rs := []string{}

	err := opt.Validate()
	if err != nil {
		return rs, err
	}

	command, err := opt.Executable()
	if err != nil {
		return rs, err
	}

	cmd := exec.Command(getShell(), "-c", command)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return rs, err
	}

	go func() {
		input(stdin)
		stdin.Close()
	}()

	stdout, err := cmd.Output()
	if err != nil {
		return rs, err
	}

	for _, key := range strings.Split(string(stdout), "\n") {
		if !empty(key) {
			rs = append(rs, strings.TrimSpace(key))
		}
	}

	return rs, nil
}

func FromArray(opt *GofiOptions, input []string) ([]string, error) {
	rs := []string{}

	err := opt.Validate()
	if err != nil {
		return rs, err
	}

	command, err := opt.Executable()
	if err != nil {
		return rs, err
	}

	cmd := exec.Command(getShell(), "-c", command)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return rs, err
	}

	go func() {
		for _, key := range input {
			fmt.Fprintln(stdin, key)
		}
		stdin.Close()
	}()

	stdout, err := cmd.Output()
	if err != nil {
		return rs, err
	}

	for _, key := range strings.Split(string(stdout), "\n") {
		if !empty(key) {
			rs = append(rs, strings.TrimSpace(key))
		}
	}

	return rs, nil
}

func FromInterface(opt *GofiOptions, input map[string]interface{}) ([]string, error) {
	rs := []string{}

	err := opt.Validate()
	if err != nil {
		return rs, err
	}

	command, err := opt.Executable()
	if err != nil {
		return rs, err
	}

	cmd := exec.Command(getShell(), "-c", command)
	cmd.Stderr = os.Stderr

	stdin, err := cmd.StdinPipe()
	if err != nil {
		return rs, err
	}

	go func() {
		temp := []string{}

		for key := range input {
			temp = append(temp, key)
		}

		sort.Strings(temp)

		for _, key := range temp {
			fmt.Fprintln(stdin, key)
		}

		stdin.Close()
	}()

	stdout, err := cmd.Output()
	if err != nil {
		return rs, err
	}

	for _, key := range strings.Split(string(stdout), "\n") {
		if !empty(key) {
			rs = append(rs, strings.TrimSpace(key))
		}
	}

	return rs, nil
}
