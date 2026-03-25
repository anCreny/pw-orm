package pworm

import (
	"testing"
)

func TestValidWhere(t *testing.T) {
	cmd := NewCommandBuilder("Test-Command").Select([]string{
		"field1", "field2", "field3",
	}...).Limit(20).Where(
		WhereCondition("field1", Equal, "5").
			AND().
			WhereCondition("field2", LessThen, "5").
			OR().
			WhereCondition("field3", NotLike, "*pow*"),
	).SetArguments(
		Argument{
			Name:  "Arg1",
			Value: 5,
		},
	).Command()

	expectCommand := "powershell -command Test-Command -Arg1 5 | Where-Object { field1 -eq '5' -and field2 -lt '5' -or field3 -notlike '*pow*' } | Select field1, field2, field3 | Select-Object -First 20 | ConvertTo-Json"
	actualCommand := cmd.String()

	if expectCommand != actualCommand {
		t.Errorf("\n commands don't match \n expect: %s \n actual: %s \n", expectCommand, actualCommand)
	}
}
