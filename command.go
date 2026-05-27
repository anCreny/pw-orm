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

func (c *Command) String() string {
	return c.command
}

func (c *Command) Cmd() *exec.Cmd {
	return exec.Command("powershell", "-command", c.command)
}

func (c *Command) Run() (result, error) {
	commandToExtract := `
	Try 
	{ 
		$res = ` + c.command + `
		@{"Output" = $res} | ConvertTo-Json -Depth 4
	} 
	Catch 
	{ 
		$res = $_  | ConvertTo-Json -Depth 3  
		@{"Error" = $res} | ConvertTo-Json
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
