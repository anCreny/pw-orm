package pworm

import "os/exec"

type Executer interface {
	Execute(string) ([]byte, error)
}

type TransientExecuter struct {
}

func (e *TransientExecuter) Execute(command string) ([]byte, error) {
	return exec.Command("powershell", "-command", command).Output()
}
