package pworm

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anCreny/pw-orm/errors"
)

type CommandBuilder struct {
	command   string
	arguments []string

	whereClause  string
	selectClause string
	limit        int
	autoConfirm  bool
	errorAction  string
	mustArray    bool

	executor Executer
}

func NewCommandBuilder(command string) *CommandBuilder {
	return &CommandBuilder{
		command:  command,
		executor: &TransientExecuter{},
	}
}

func (c *CommandBuilder) Build() (*Command, error) {

	// Собираем команду

	command := c.command

	if len(c.arguments) != 0 {
		argsString := strings.Join(c.arguments, " ")
		command = fmt.Sprintf("%s %s", command, argsString)
	}

	if c.autoConfirm {
		command = fmt.Sprintf("%s -Confirm:$false", command)
	}

	if c.errorAction != "" {
		command = fmt.Sprintf("%s -EA %s", command, c.errorAction)
	}

	if c.whereClause != "" {
		whereString := fmt.Sprintf("Where-Object {%s}", c.whereClause)

		command = fmt.Sprintf("%s | %s", command, whereString)
	}

	if c.selectClause != "" {
		selectString := fmt.Sprintf("| Select %s", c.selectClause)

		command = fmt.Sprintf("%s %s", command, selectString)
	}

	if c.limit != 0 {

		limitString := fmt.Sprintf("| Select-Object -First %d", c.limit)

		command = fmt.Sprintf("%s %s", command, limitString)
	}

	if c.mustArray {
		command = fmt.Sprintf("@(%s)", command)
	}

	// Проверив команду на валидность с помощью встроеного PowerShell механизма

	validateCommand := `
	Try {
		$userCommand = @'
		` + command + `
'@

		$command = [ScriptBlock]::Create($userCommand)
	}
	Catch {
		@{"Error" = $_} | ConvertTo-Json -Depth 3
	}

	`

	res, err := c.executor.Execute(validateCommand)
	if err != nil {
		return nil, err
	}

	if res != nil {
		var result ValudateResult

		if err := json.Unmarshal(res, &result); err != nil {
			return nil, fmt.Errorf("ошибка декодирования результата проверки команды: %v", err)
		}

		if result.Error != nil {
			return nil, fmt.Errorf("ошибка синтаксиса команды: %s", result.Error.ErrorDetails.Message)
		}
	}

	return &Command{
		command:  command,
		executor: c.executor,
	}, nil
}

type ValudateResult struct {
	Error *errors.Error `json:"Error"`
}

func (c *CommandBuilder) AutoConfirm() *CommandBuilder {

	c.autoConfirm = true

	return c
}

func (c *CommandBuilder) MustArray() *CommandBuilder {
	c.mustArray = true
	return c
}
