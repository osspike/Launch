package terminal

import (
	"os"

	"golang.org/x/term"
)

type Console struct {
	oldState *term.State
}

func NewConsole() *Console {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	if err != nil {
		panic(err)
	}
	return &Console{oldState: oldState}
}

func (c *Console) Restore() {
	term.Restore(int(os.Stdin.Fd()), c.oldState)
}

func (c *Console) Write(b []byte) (int, error) {
	return os.Stdout.Write(b)
}

func (c *Console) Read(b []byte) (int, error) {
	return os.Stdin.Read(b)
}
