package pworm

import "fmt"

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
