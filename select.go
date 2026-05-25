package pworm

// if As == "" { As = Name }
//
// if $_.'Name' == NotFound { result.Out.'Name' = nil }
type Field struct {
	Name string
	As   string
}

func (c *CommandBuilder) Select(fields ...Field) *CommandBuilder {
	c.selectFields = append(c.selectFields, fields...)

	return c
}
