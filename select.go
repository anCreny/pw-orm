package pworm

import "fmt"

// if As == "" { As = Name }
//
// if $_.'Name' == NotFound { result.Out.'Name' = nil }
type Field struct {
	Name string
	As   string
}

func (c *CommandBuilder) Select(fields ...Field) *CommandBuilder {

	for _, field := range fields {
		if field.As == "" {
			field.As = field.Name
		}

		if c.selectClause == "" {
			c.selectClause = fmt.Sprintf("@{Name='%s'; Expression={$_.%s}}", field.As, field.Name)
			continue
		}

		c.selectClause = fmt.Sprintf("%s, @{Name='%s'; Expression={$_.%s}}", c.selectClause, field.As, field.Name)
	}

	return c
}
