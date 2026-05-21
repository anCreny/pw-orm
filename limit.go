package pworm

func (c *CommandBuilder) Limit(count int) *CommandBuilder {
	if count < 0 {
		count = 0
	}

	c.limit = count

	return c
}
