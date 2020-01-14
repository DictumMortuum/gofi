package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"log"
)

func isExecutable(name string) bool {
	_, err := exec.LookPath(name)
	return err == nil
}

func startedFromTerminal() bool {
	return os.Getenv("TERM") != ""
}

func getShell() string {
	shell := os.Getenv("SHELL")

	if len(shell) == 0 {
		shell = "sh"
	}

	return shell
}

func withFilter(command string, input func(in io.WriteCloser)) string {
	cmd := exec.Command(getShell(), "-c", command)
	cmd.Stderr = os.Stderr
	in, _ := cmd.StdinPipe()

	go func() {
		input(in)
		in.Close()
	}()

	result, _ := cmd.Output()
	temp := string(result)
	temp = strings.Replace(temp, "\n", " ", -1)
	temp = strings.TrimSpace(temp)

	return temp
}

func executable() string {
	if startedFromTerminal() && isExecutable("fzf") {
		return "fzf"
	}

	if isExecutable("rofi") {
		return "rofi"
	} else if isExecutable("dmenu") {
		return "dmenu"
	}

	log.Fatal("No executable found")
	return ""
}

func plain(desc string) string {
	switch executable() {
	case "rofi":
		return fmt.Sprintf("rofi -dmenu -multi-select -matching fuzzy -i -p '%s'", desc)
	case "dmenu":
		return fmt.Sprintf("dmenu")
	case "fzf":
		return fmt.Sprintf("fzf -m -i")
	default:
		return "echo"
	}
}

func Plain(desc string, input func(in io.WriteCloser)) string {
	return withFilter(plain(desc), input)
}

func multi(desc string) string {
	switch executable() {
	case "rofi":
		return fmt.Sprintf("rofi -dmenu -multi-select -matching fuzzy -i -p '%s'", desc)
	case "dmenu":
		return fmt.Sprintf("dmenu")
	case "fzf":
		return fmt.Sprintf("fzf -m -i")
	default:
		return "echo"
	}
}

func Multi(desc string, input func(in io.WriteCloser)) string {
	return withFilter(multi(desc), input)
}

func main() {
	Plain("env", func(in io.WriteCloser) {
		files := []string{"a", "b", "c"}

		for _, file := range files {
			fmt.Fprintln(in, file)
		}
	})
}