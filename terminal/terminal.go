package terminal

import (
	"io"
	"launch/color"
	"launch/config"
	"os"
	"os/exec"

	"golang.org/x/term"
)

type Terminal struct {
	console  *Console
	terminal *term.Terminal
	running  bool
}

func NewTerminal() *Terminal {
	c := NewConsole()
	return &Terminal{console: c, terminal: term.NewTerminal(c, "> "), running: true}
}

func (t *Terminal) Close() {
	t.running = false
	t.console.Restore()
}

func (t *Terminal) ProcessInput() (done chan struct{}) {
	done = make(chan struct{})
	go func() {
		for t.running {
			command, err := t.terminal.ReadLine()
			if err == nil {
				run("cmd", "/c", command)
			} else if err == io.EOF {
				done <- struct{}{}
			}
		}
	}()

	return
}

func (t *Terminal) StartProcess(p config.Process) (*exec.Cmd, error) {
	cmd := exec.Command(p.Path, p.Args...)

	out := &color.Color{
		Value: p.Color.Value,
		Out:   t.terminal,
	}

	cmd.Stdout = out
	cmd.Stderr = out
	if p.Cwd != "" {
		cmd.Dir = p.Cwd
	}
	err := cmd.Start()
	return cmd, err
}

func run(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
