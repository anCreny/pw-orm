package pworm

import (
	"fmt"
	"strings"
)

type CommandBuilder struct {
	command   string
	arguments []string

	whereClause  string
	selectFields []Field
	limit        int
	autoConfirm  bool
	errorAction  string

	executor Executer
}

func NewCommandBuilder(command string) *CommandBuilder {
	return &CommandBuilder{
		command:  command,
		executor: &TransientExecuter{},
	}
}

func (c *CommandBuilder) Build() *Command {

	command := c.command

	if len(c.arguments) != 0 {
		argsString := strings.Join(c.arguments, " ")
		command = fmt.Sprintf("%s %s", command, argsString)
	}

	if c.autoConfirm {
		command = fmt.Sprintf("%s -Confirm:$false", command)
	}

	command = fmt.Sprintf("%s -EA Stop", command)

	if c.whereClause != "" {
		whereString := fmt.Sprintf("Where-Object {%s}", c.whereClause)

		command = fmt.Sprintf("%s | %s", command, whereString)
	}

	if len(c.selectFields) != 0 {
		var fields []string
		for _, selectField := range c.selectFields {
			// Если имя пустое, пропускаем
			if selectField.Name == "" {
				continue
			}

			// Если алиас не указан, устанавливаем
			// его в качестве имени
			if selectField.As == "" {
				selectField.As = selectField.Name
			}

			field := fmt.Sprintf("@{Name='%s'; Expression={$_.%s}}", selectField.As, selectField.Name)

			fields = append(fields, field)
		}
		selectString := strings.Join(fields, ", ")

		selectString = fmt.Sprintf("| Select %s", selectString)

		command = fmt.Sprintf("%s %s", command, selectString)
	}

	if c.limit != 0 {

		limitString := fmt.Sprintf("| Select-Object -First %d", c.limit)

		command = fmt.Sprintf("%s %s", command, limitString)
	}

	return &Command{
		command:  command,
		executor: c.executor,
	}
}

func (c *CommandBuilder) AutoConfirm() *CommandBuilder {

	c.autoConfirm = true

	return c
}
