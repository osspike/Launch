package color

import (
	"encoding/json"
	"io"
	"os"
	"strings"
)

type Color struct {
	Value byte
	Out   io.Writer
}

const (
	GREY byte = 30 + iota
	RED
	GREEN
	YELLOW
	BLUE
	MAGENTA
	CYAN
	WHITE
)

func (c *Color) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return err
	}

	var v byte
	switch strings.ToLower(s) {
	case "grey":
		v = GREY
	case "red":
		v = RED
	case "green":
		v = GREEN
	case "yellow":
		v = YELLOW
	case "blue":
		v = BLUE
	case "magenta":
		v = MAGENTA
	case "cyan":
		v = CYAN
	default:
		v = WHITE
	}
	*c = Color{Value: v}

	return nil
}

func (c *Color) Reset(other Color) {
	*c = other
}

func (c *Color) Write(data []byte) (int, error) {
	p := prefix(*c)
	s := suffix()

	message := make([]byte, 0, len(p)+len(data)+len(s))
	message = append(message, p...)
	message = append(message, data...)
	message = append(message, s...)

	out := c.Out
	if out == nil {
		out = os.Stdout
	}

	n, err := out.Write(message)
	if err != nil {
		return n, err
	}

	return len(data), nil
}

func prefix(c Color) []byte {
	b := make([]byte, 5)
	b[0] = '\033'
	b[1] = '['
	b[2] = byte(c.Value)/10%10 + 48 // to ASCII
	b[3] = byte(c.Value)%10 + 48
	b[4] = 'm'
	return b
}

func suffix() []byte {
	return []byte("\033[0m")
}
