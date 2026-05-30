package pworm

import (
	"github.com/anCreny/pw-orm/sessions"
)

type Operator struct {
	s *sessions.Session
}

func (o *Operator) NewCommandBuilder(command string) *SCommandBuilder {
	return &SCommandBuilder{
		c: &CommandBuilder{
			command:  command,
			executor: o.s,
		},
		o: o,
	}
}

func (o *Operator) RawCommand(rawCommand string) *SCommand {
	return &SCommand{
		c: &Command{
			command:  rawCommand,
			executor: o.s,
		},
		o: o,
	}
}

// SessionCommandBuilder
type SCommandBuilder struct {
	c *CommandBuilder
	o *Operator
}

func (s *SCommandBuilder) Where(cond Condition) *SCommandBuilder {
	s.c = s.c.Where(cond)
	return s
}

func (s *SCommandBuilder) Limit(count int) *SCommandBuilder {
	s.c = s.c.Limit(count)
	return s
}

func (s *SCommandBuilder) Select(fields ...Field) *SCommandBuilder {
	s.c = s.c.Select(fields...)
	return s
}

func (s *SCommandBuilder) WithArguments(args ...Argument) *SCommandBuilder {
	s.c = s.c.WithArguments(args...)
	return s
}

func (s *SCommandBuilder) Build() (*SCommand, error) {
	command, err := s.c.Build()
	return &SCommand{
		c: command,
		o: s.o,
	}, err
}

type SCommand struct {
	c *Command
	o *Operator
}

func (c *SCommand) Run() (result, error) {
	return c.c.Run()
}
