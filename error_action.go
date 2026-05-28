package pworm

type errorAction string

const (
	EA_STOP       errorAction = "Stop"
	EA_CONTINUE   errorAction = "Continue"
	EA_IGNORE     errorAction = "Ignore"
	EA_SILENT     errorAction = "SilentlyContinue"
	EA_PREFERENCE errorAction = "ActionPreference"
)

func (c *CommandBuilder) ErrorAction(action errorAction) *CommandBuilder {
	c.errorAction = string(action)
	return c
}
