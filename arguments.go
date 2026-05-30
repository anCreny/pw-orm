package pworm

import "fmt"

func (c *CommandBuilder) WithArguments(args ...Argument) *CommandBuilder {
	arguments := make([]Arg, 0, len(args))
	for _, arg := range args {
		arguments = append(arguments, arg.toArgument())
	}

	c.arguments = arguments

	return c
}

type Arg struct {
	Value string
	Row   string
}

type Argument interface {
	toArgument() Arg
}

type StringArg struct {
	Name  string
	Value string
}

func (a *StringArg) toArgument() Arg {
	return Arg{
		Value: fmt.Sprintf("-%s \"%s\"", a.Name, a.Value),
		Row:   fmt.Sprintf("\"%s\" = \"%s\"", a.Name, a.Value),
	}
}

type IntArg struct {
	Name  string
	Value int
}

func (a *IntArg) toArgument() Arg {
	return Arg{
		Value: fmt.Sprintf("-%s %d", a.Name, a.Value),
		Row:   fmt.Sprintf("\"%s\" = %d", a.Name, a.Value),
	}
}

type ConstantArg struct {
	Name  string
	Value string
}

func (a *ConstantArg) toArgument() Arg {
	return Arg{
		Value: fmt.Sprintf("-%s %s", a.Name, a.Value),
		Row:   fmt.Sprintf("\"%s\" = \"%s\"", a.Name, a.Value),
	}
}

type NilArg struct {
	Name string
}

func (a *NilArg) toArgument() Arg {
	return Arg{
		Value: fmt.Sprintf("-%s", a.Name),
		Row:   fmt.Sprintf("\"%s\" = $true", a.Name),
	}
}

type SVariableArg struct {
	Name  string
	Value SVariable
}

func (a *SVariableArg) toArgument() Arg {
	return Arg{
		Value: fmt.Sprintf("-%s %s", a.Name, a.Value.PW()),
		Row:   fmt.Sprintf("\"%s\" = %s", a.Name, a.Value.PW()),
	}
}
