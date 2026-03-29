package pworm

import (
	"fmt"
	"os/exec"
	"strings"
	"testing"
)

type CommandBuilder struct {
	command   string
	arguments []string

	whereClause  string
	selectFields []string
	limit        int
	autoConfirm  bool
}

func NewCommandBuilder(command string) *CommandBuilder {
	return &CommandBuilder{
		command: command,
	}
}

// valueConstant - тип данных, обозначающий константу.
//
// Нужен для передачи значения аргумента команды
// как константы, без кавычек ("")
type valueConstant string

func Constant(v string) valueConstant {
	return valueConstant(v)
}

// valueNil - тип данных, обозначающий пустое значение.
//
// Нужен для передачи пустого значения, когда
// аргумент не требует передачи никакого значения.
type valueNil struct{}

func Nil() valueNil {
	return valueNil{}
}

type Argument struct {
	Name  string
	Value any
}

func (a *Argument) toArgument() string {

	switch any(a.Value).(type) {
	case int, int32, int64:
		return fmt.Sprintf("-%s %d", a.Name, a.Value.(int))
	case string:
		return fmt.Sprintf("-%s \"%s\"", a.Name, a.Value.(string))
	case valueConstant:
		return fmt.Sprintf("-%s %s", a.Name, a.Value.(string))
	case valueNil:
		return fmt.Sprintf("-%s", a.Name)
	default:
		panic("неподдерживаемый тип аргумента")
	}
}

func (c *CommandBuilder) SetArguments(args ...Argument) *CommandBuilder {
	arguments := make([]string, 0, len(args))
	for _, arg := range args {
		arguments = append(arguments, arg.toArgument())
	}

	c.arguments = arguments

	return c
}

type Condition string

var (
	Equal    Condition = "-eq"
	NotEqual Condition = "-ne"
	LessThen Condition = "-lt"
	Like     Condition = "-like"
	NotLike  Condition = "-notlike"
	RegExp   Condition = "-match"
	Contains Condition = "-contains"
)

type whereClause struct {
	clause string
}

type connector struct {
	w *whereClause
}

func (w *whereClause) OR() *connector {
	w.clause = fmt.Sprintf("%s-or ", w.clause)

	return &connector{
		w: w,
	}
}

func (w *whereClause) AND() *connector {
	w.clause = fmt.Sprintf("%s-and ", w.clause)

	return &connector{
		w: w,
	}
}

func (c *connector) WhereCondition(field string, cond Condition, value string) *whereClause {
	w := c.w

	w.clause = fmt.Sprintf("%s%s %s '%s' ", w.clause, field, cond, value)

	return w
}

func WhereCondition(field string, cond Condition, value string) *whereClause {
	w := &whereClause{}

	w.clause = fmt.Sprintf(" %s %s '%s' ", field, cond, value)

	return w
}

func (c *CommandBuilder) Where(w *whereClause) *CommandBuilder {

	c.whereClause = w.clause

	return c
}

func (c *CommandBuilder) Select(fields ...string) *CommandBuilder {
	m := make([]string, 0, len(fields))

	m = append(m, fields...)

	c.selectFields = append(c.selectFields, m...)

	return c
}

func (c *CommandBuilder) Limit(count int) *CommandBuilder {
	if count < 0 {
		count = 0
	}

	c.limit = count

	return c
}

func (c *CommandBuilder) Command() *exec.Cmd {

	command := c.command

	if len(c.arguments) != 0 {
		argsString := strings.Join(c.arguments, " ")
		command = fmt.Sprintf("%s %s", command, argsString)
	}

	if c.autoConfirm {
		command = fmt.Sprintf("%s -Confirm:$false", command)
	}

	if c.whereClause != "" {
		whereString := fmt.Sprintf("Where-Object {%s}", c.whereClause)

		command = fmt.Sprintf("%s | %s", command, whereString)
	}

	if len(c.selectFields) != 0 {
		selectString := strings.Join(c.selectFields, ", ")

		selectString = fmt.Sprintf("| Select %s", selectString)

		command = fmt.Sprintf("%s %s", command, selectString)
	}

	if c.limit != 0 {

		limitString := fmt.Sprintf("| Select-Object -First %d", c.limit)

		command = fmt.Sprintf("%s %s", command, limitString)
	}

	command = fmt.Sprintf("%s | ConvertTo-Json", command)

	return exec.Command("powershell", "-command", command)
}

func (c *CommandBuilder) AutoConfirm() *CommandBuilder {

	c.autoConfirm = true

	return c
}

func TestValidWherer(t *testing.T) {
	cmd := NewCommandBuilder("Test-Command").Select([]string{
		"field1", "field2", "field3",
	}...).Limit(20).Where(
		WhereCondition("field1", Equal, "5").
			AND().
			WhereCondition("field2", LessThen, "5"),
	).SetArguments(
		Argument{
			"Arg1",
			5,
		},
	).Command()

	expectCommand := "powershell -command Test-Command -Arg1 5 | Where-Object { field1 -eq '5' -and field2 -lt '5' } | Select field1, field2, field3 | Select-Object -First 20 | ConvertTo-Json"
	actualCommand := cmd.String()

	if expectCommand != actualCommand {
		t.Errorf("\n commands don't match \n expect: %s \n actual: %s \n", expectCommand, actualCommand)
	}
}
