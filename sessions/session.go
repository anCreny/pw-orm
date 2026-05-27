package sessions

import (
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"

	"github.com/anCreny/pw-orm/helpers"
)

type Session struct {
	cmd *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	stdoutClose chan struct{}
	stderrClone chan struct{}

	ID string // Идентификатор сессии для управления scope в powershell
}

func Start() (*Session, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "-")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	s := &Session{
		ID:     helpers.GenerateRandomString(10),
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}

	if _, err := s.Execute(fmt.Sprintf("$global:%s = @{}", s.ID)); err != nil {
		return nil, err
	}

	return s, nil
}

func (s *Session) Execute(command string) ([]byte, error) {
	wrappedCommand := `
		$stdoutStr = & { ` + command + ` }

    # ПРЕДОХРАНИТЕЛЬ: Если команда ничего не вывела (например, это присвоение переменной),
    # мы принудительно устанавливаем валидный пустой JSON, чтобы Go не завис.
    if ([string]::IsNullOrWhiteSpace($stdoutStr)) {
        $stdoutStr = "{}"
    }

    # Перехватываем системные ошибки (если они были)
    $stderrStr = $Error | Out-String; $Error.Clear()

    # Переводим данные в сырые байты UTF-8
    $outBytes = [System.Text.Encoding]::UTF8.GetBytes($stdoutStr.Trim())
    $errBytes = [System.Text.Encoding]::UTF8.GetBytes($stderrStr.Trim())

    function Send-Packet([byte[]]$bytes, $isError) {
        $stream = if ($isError) { [System.Console]::OpenStandardError() } else { [System.Console]::OpenStandardOutput() }
        $lenBytes = [System.BitConverter]::GetBytes($bytes.Length)
        $stream.Write($lenBytes, 0, 4)
        if ($bytes.Length -gt 0) { $stream.Write($bytes, 0, $bytes.Length) }
    }

    # Отправляем пакеты
    Send-Packet $outBytes $false
    Send-Packet $errBytes $true
	`

	if _, err := fmt.Fprintln(s.stdin, wrappedCommand); err != nil {
		return nil, err
	}

	// Читаем оба потока параллельно, чтобы избежать взаимной блокировки (Deadlock) буферов ОС
	type packetResult struct {
		data []byte
		err  error
	}
	stdoutChan := make(chan packetResult, 1)
	stderrChan := make(chan packetResult, 1)

	go func() {
		data, err := readPacket(s.stdout)
		stdoutChan <- packetResult{data, err}
	}()

	go func() {
		data, err := readPacket(s.stderr)
		stderrChan <- packetResult{data, err}
	}()

	resStdout := <-stdoutChan
	resStderr := <-stderrChan

	// Если произошел сбой самого системного пайпа — возвращаем ошибку Go
	if resStdout.err != nil {
		return nil, resStdout.err
	}
	if resStderr.err != nil {
		return nil, resStderr.err
	}

	// ЛОГИКА ВЫБОРА ВЫВОДА:
	// Если стандартный вывод пустой (или равен nil/0 байт), а в ошибках что-то есть
	if len(resStdout.data) == 0 && len(resStderr.data) > 0 {
		return resStderr.data, nil
	}

	// Во всех остальных случаях (есть только stdout, или есть оба) возвращаем stdout
	return resStdout.data, nil
}

func readPacket(r io.Reader) ([]byte, error) {
	var length int32
	if err := binary.Read(r, binary.LittleEndian, &length); err != nil {
		return nil, err
	}

	if length == 0 {
		return []byte{}, nil
	}

	buf := make([]byte, length)
	if _, err := io.ReadFull(r, buf); err != nil {
		return nil, err
	}

	return buf, nil
}

func (s *Session) Close() {
	fmt.Fprintln(s.stdin, "exit")
	s.stdin.Close()
	_ = s.cmd.Wait()
}
