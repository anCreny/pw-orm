package pworm

import (
	"fmt"
	"strings"
)

type op string

var (
	Equal    op = "-eq"
	NotEqual op = "-ne"
	LessThen op = "-lt"
	Like     op = "-like"
	NotLike  op = "-notlike"
	RegExp   op = "-match"
	Contains op = "-contains"
)

func (c *CommandBuilder) Where(cond Condition) *CommandBuilder {

	c.whereClause = string(cond)

	return c
}

type Condition string

func (c Condition) ToString() string {
	return string(c)
}

func Clause(field string, o op, value string) Condition {
	return Condition(fmt.Sprintf(" $_.%s %s '%s' ", field, o, value))
}

func ANDClause(conds ...Condition) Condition {
	res := joinConditions(conds, "-and")

	return Condition(fmt.Sprintf(" (%s) ", res))
}

func ORClause(conds ...Condition) Condition {
	res := joinConditions(conds, "-or")

	return Condition(fmt.Sprintf(" (%s) ", res))
}

func joinConditions(conds []Condition, connector string) string {
	var res []string

	for _, cond := range conds {
		res = append(res, cond.ToString())
	}

	return strings.Join(res, connector)
}
