package pworm

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/anCreny/pw-orm/errors"
)

type CommandBuilder struct {
	command   string
	arguments []Arg

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
		var argsString []string

		for _, arg := range c.arguments {
			argsString = append(argsString, arg.Value)
		}

		args := strings.Join(argsString, " ")
		command = fmt.Sprintf("%s %s", command, args)
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

	var commandPartCheckParams string

	if len(c.arguments) > 0 {

		var argRows []string

		for _, arg := range c.arguments {
			argRows = append(argRows, arg.Row)
		}

		args := strings.Join(argRows, "\n")

		commandPartCheckParams = `
		$parameters = @{
			` + args + `
			}

			$cmdInfo = Get-Command -Name "` + c.command + `" -ErrorAction Stop

      foreach ($key in $parameters.Keys) {
        if (-not $cmdInfo.Parameters.ContainsKey($key)) {
             throw [System.Management.Automation.ParameterBindingException]"Parameter '$key' not found for the command."
        }
     	}
		`

	}

	validateCommand := `
	Try {
			$userCommand = @'
			` + command + `
'@

			$null = [ScriptBlock]::Create($userCommand)

			` + commandPartCheckParams + `

	}
	Catch {
		 	if ($_.Exception -and $_.Exception.Data) { $_.Exception.Data.Clear() }
     	@{"Error" = $_ | Select-Object Exception, CategoryInfo, ErrorDetails} | ConvertTo-Json -Depth 3
	}
	`

	res, err := c.executor.Execute(validateCommand)
	if err != nil {
		return nil, err
	}

	res = bytes.TrimSpace(res)

	if res != nil {
		var result ValudateResult

		if err := json.Unmarshal(res, &result); err != nil {
			return nil, fmt.Errorf("ошибка декодирования результата проверки команды: %v", err)
		}

		if result.Error != nil {
			return nil, fmt.Errorf("ошибка синтаксиса команды: %s", result.Error.Exception.Message)
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
