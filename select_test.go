package pworm

import "testing"

func TestValidSelect1(t *testing.T) {
	c := NewCommandBuilder("select")
	c.Select(Field{Name: "id", As: "ID"})
	c.Select(Field{Name: "name", As: "Name"})
	c.executor = &FakeExecuter{}
	command, err := c.Build()
	if err != nil {
		t.Fatal(err)
	}

	currentCommand := command.ToString()
	exprectedCommand := "select | Select @{Name='ID'; Expression={$_.id}}, @{Name='Name'; Expression={$_.name}}"

	if currentCommand != exprectedCommand {
		t.Errorf("\n commands don't match \n expect: %s \n actual: %s \n", exprectedCommand, currentCommand)
	}

}
