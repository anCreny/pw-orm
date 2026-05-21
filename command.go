package pworm

import (
	"encoding/json"
	"fmt"
	"os/exec"
)

type Command struct {
	command string
}

func (c *Command) String() string {
	return c.command
}

func (c *Command) Cmd() *exec.Cmd {
	return exec.Command("powershell", "-command", c.command)
}

func (c *Command) Run() (result, error) {
	commandToExtract := `Try 
	{ $res = ` + c.command + ` ; $res = '{"Output": ; + $res + "}" ; Write-Host $res } 
	Catch 
	{ $res = $_  | ConvertTo-Json ; $res = '{"Error": ' + $res + "}" ; Write-Host $res }`

	commandOut, err := exec.Command("powershell", "-command", commandToExtract).Output()
	if err != nil {
		return result{}, fmt.Errorf("error on start running command(%s): %s", c.command, err)
	}

	var r result

	if err = json.Unmarshal(commandOut, &r); err != nil {
		return result{}, fmt.Errorf("error on unmarshal command out(%s): %s", string(commandOut), err)
	}

	return r, nil
}
