package lib

import (
	"errors"
	"github.com/frobware/go-wordexp"
	"os/exec"
)

type Command struct {
	Path string
	Args []string
}

func (c *Command) UnmarshalText(text []byte) error {
	argv, err := wordexp.Expand(string(text))
	if err != nil {
		return err
	}
	if len(argv) == 0 {
		return errors.New("command string is too short")
	}
	c.Path = argv[0]
	c.Args = argv[1:]
	return nil
}

func (c *Command) GetCmd() *exec.Cmd {
	return exec.Command(c.Path, c.Args...)
}
