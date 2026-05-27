package pworm

// RawCommand превращает строку с командой
// для PowerShell в структуру Command для
// возможности взаимодействия с командой через
// pworm.
//
// pworm НЕ вмешивается в синтаксис переданной
// команды. За все ошибки ее выполнения отвечает
// тот, кто ее передал.
func RawCommand(pwCommand string) *Command {
	return &Command{command: pwCommand}
}
