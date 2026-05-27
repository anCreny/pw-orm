package sessions

import (
	"encoding/binary"
	"fmt"
	"io"
	"os/exec"
)

type Session struct {
	cmd *exec.Cmd

	stdin  io.WriteCloser
	stdout io.ReadCloser
	stderr io.ReadCloser

	stdoutClose chan struct{}
	stderrClone chan struct{}
}

func Start() (*Session, error) {
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", "-")

	stdin, _ := cmd.StdinPipe()
	stdout, _ := cmd.StdoutPipe()
	stderr, _ := cmd.StderrPipe()

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	return &Session{
		cmd:    cmd,
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
	}, nil
}

func (s *Session) Execute(command string) ([]byte, error) {
	wrappedCommand := `
		$stdoutStr = & { ` + command + ` } | Out-String
		$stderrStr = $Error | Out-String; $Error.Clear()

		$outBytes = [System.Text.Encoding]::UTF8.GetBytes($stdoutStr.Trim())
		$errBytes = [System.Text.Encoding]::UTF8.GetBytes($stderrStr.Trim())

		function Send-Packet([byte[]]$bytes) {
			$stream = [System.Console]::OpenStandardOutput()
			$lenBytes = [System.BitConverter]::GetBytes($bytes.Length)
			$stream.Write($lenBytes, 0, 4)
			if ($bytes.Length -gt 0) { $stream.Write($bytes, 0, $bytes.Length) }
		}

		# Всегда отправляем ровно 2 пакета в фиксированном порядке
		Send-Packet $outBytes
		
		# Для Stderr пишем напрямую в поток ошибок Windows
		$errStream = [System.Console]::OpenStandardError()
		$errLenBytes = [System.BitConverter]::GetBytes($errBytes.Length)
		$errStream.Write($errLenBytes, 0, 4)
		if ($errBytes.Length -gt 0) { $errStream.Write($errBytes, 0, $errBytes.Length) }
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
