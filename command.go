package pworm

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Command struct {
	command  string
	executor Executer
}

func (c *Command) ToString() string {
	return c.command
}

func (c *Command) Cmd() *exec.Cmd {
	return exec.Command("powershell", "-command", c.command)
}

func (c *Command) Run() (result, error) {
	// Если команда выполеняется успешно,
	// создастатся глобальная переменная $global:Output.
	// Для того, чтобы избежать проблем, связанных с чтением
	// старой глобальной переменной $global:Output во время
	// ошибки выполнения следующей команды, удалим ее.
	commandToExtract := `
	Try 
	{ 

		$res = ` + c.command + `

		$global:Output_` + scopeID + ` = $res

		if ($null -eq $res) {
         '{"Output": null}'
    } else {
         @{"Output" = $res} | ConvertTo-Json -Depth 4
    }
	} 
	Catch 
	{ 
		Remove-Variable -Name "Output_` + scopeID + `" -Scope Global -ErrorAction SilentlyContinue
		@{"Error" = $_} | ConvertTo-Json -Depth 3
	}
	`

	commandOut, err := c.executor.Execute(commandToExtract)
	if err != nil {
		return result{}, fmt.Errorf("error on run command(%s): %s", c.command, err)
	}

	var r result

	if err = json.Unmarshal(commandOut, &r); err != nil {
		return result{}, fmt.Errorf("error on unmarshal command out(%s): %s", string(commandOut), err)
	}

	return r, nil
}
