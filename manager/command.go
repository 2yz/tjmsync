package manager

import (
	"errors"
	"strings"
)

type Command struct {
	Path string
	Args []string
}

func ParseCommand(str string) Command {
	c := Command{}
	c.UnmarshalText([]byte(str))
	return c
}

func (c *Command) UnmarshalText(text []byte) (err error) {
	argv := strings.Fields(string(text))
	if err != nil {
		return
	}
	if len(argv) == 0 {
		err = errors.New("command string is too short")
		return
	}
	c.Path = argv[0]
	c.Args = argv[1:]
	return
}

func (c *Command) MarshalText() ([]byte, error) {
	str := c.Path + " " + strings.Join(c.Args, " ")
	return []byte(str), nil
}

func (c *Command) GetCmd() []string {
	return append([]string{c.Path}, c.Args...)
}
