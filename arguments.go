package pworm

import "fmt"

func (c *CommandBuilder) WithArguments(args ...Argument) *CommandBuilder {
	arguments := make([]string, 0, len(args))
	for _, arg := range args {
		arguments = append(arguments, arg.toArgument())
	}

	c.arguments = arguments

	return c
}

type Argument interface {
	toArgument() string
}

type StringArg struct {
	Name  string
	Value string
}

func (a *StringArg) toArgument() string {
	return fmt.Sprintf("-%s \"%s\"", a.Name, a.Value)
}

type IntArg struct {
	Name  string
	Value int
}

func (a *IntArg) toArgument() string {
	return fmt.Sprintf("-%s %d", a.Name, a.Value)
}

type ConstantArg struct {
	Name  string
	Value string
}

func (a *ConstantArg) toArgument() string {
	return fmt.Sprintf("-%s %s", a.Name, a.Value)
}

type NilArg struct {
	Name string
}

func (a *NilArg) toArgument() string {
	return fmt.Sprintf("-%s", a.Name)
}

type SVariableArg struct {
	Name  string
	Value SVariable
}

func (a *SVariableArg) toArgument() string {
	return fmt.Sprintf("-%s %s", a.Name, a.Value.PW())
}
