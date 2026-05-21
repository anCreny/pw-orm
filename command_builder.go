package pworm

import (
	"fmt"
	"strings"
)

type CommandBuilder struct {
	command   string
	arguments []string

	whereClause  string
	selectFields []string
	limit        int
	autoConfirm  bool
}

func NewCommandBuilder(command string) *CommandBuilder {
	return &CommandBuilder{
		command: command,
	}
}

func (c *CommandBuilder) Command() *Command {

	command := c.command

	if len(c.arguments) != 0 {
		argsString := strings.Join(c.arguments, " ")
		command = fmt.Sprintf("%s %s", command, argsString)
	}

	if c.autoConfirm {
		command = fmt.Sprintf("%s -Confirm:$false", command)
	}

	if c.whereClause != "" {
		whereString := fmt.Sprintf("Where-Object {%s}", c.whereClause)

		command = fmt.Sprintf("%s | %s", command, whereString)
	}

	if len(c.selectFields) != 0 {
		selectString := strings.Join(c.selectFields, ", ")

		selectString = fmt.Sprintf("| Select %s", selectString)

		command = fmt.Sprintf("%s %s", command, selectString)
	}

	if c.limit != 0 {

		limitString := fmt.Sprintf("| Select-Object -First %d", c.limit)

		command = fmt.Sprintf("%s %s", command, limitString)
	}

	command = fmt.Sprintf("%s | ConvertTo-Json", command)

	return &Command{
		command: command,
	}
}

func (c *CommandBuilder) AutoConfirm() *CommandBuilder {

	c.autoConfirm = true

	return c
}
