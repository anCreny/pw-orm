package pworm

func (c *CommandBuilder) Select(fields ...string) *CommandBuilder {
	c.selectFields = append(c.selectFields, fields...)

	return c
}
