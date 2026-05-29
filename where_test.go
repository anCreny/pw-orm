package pworm

import (
	"testing"
)

// field1 = 5 AND field2 < 5
func TestValidWhere1(t *testing.T) {
	command := NewCommandBuilder("Test-Command").
		Where(
			ANDClause(
				Clause("field1", Equal, "5"),
				Clause("field2", LessThen, "5"),
			),
		).WithArguments(
		&IntArg{
			"Arg1",
			5,
		},
	).Build()

	expectCommand := "Test-Command -Arg1 5 | Where-Object { ( $_.field1 -eq '5' -and $_.field2 -lt '5' ) }"
	actualCommand := command.String()

	if expectCommand != actualCommand {
		t.Errorf("\n commands don't match \n expect: %s \n actual: %s \n", expectCommand, actualCommand)
	}
}

// (A = 1 or A = 2) and (B = 3 or B = 4)
func TestValidWhere2(t *testing.T) {
	command := NewCommandBuilder("Test-Command").
		Where(
			ANDClause(
				ORClause(
					Clause("A", Equal, "1"),
					Clause("A", Equal, "2"),
				),
				ORClause(
					Clause("B", Equal, "3"),
					Clause("B", Equal, "4"),
				),
			),
		).WithArguments(
		&IntArg{
			"Arg1",
			5,
		},
	).Build()

	expectCommand := "Test-Command -Arg1 5 | Where-Object { ( ( $_.A -eq '1' -or $_.A -eq '2' ) -and ( $_.B -eq '3' -or $_.B -eq '4' ) ) }"
	actualCommand := command.String()

	if expectCommand != actualCommand {
		t.Errorf("\n commands don't match \n expect: %s \n actual: %s \n", expectCommand, actualCommand)
	}
}
