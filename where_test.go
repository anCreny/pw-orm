package pworm

import (
	"testing"
)

func TestValidWhere(t *testing.T) {
	command := NewCommandBuilder("Test-Command").
		Where(
			WhereCondition("field1", Equal, "5").
				AND().
				WhereCondition("field2", LessThen, "5"),
		).WithArguments(
		&IntArg{
			"Arg1",
			5,
		},
	).Build()

	expectCommand := "Test-Command -Arg1 5 | Where-Object { $_.field1 -eq '5' -and $_.field2 -lt '5' } | ConvertTo-Json"
	actualCommand := command.String()

	if expectCommand != actualCommand {
		t.Errorf("\n commands don't match \n expect: %s \n actual: %s \n", expectCommand, actualCommand)
	}
}
