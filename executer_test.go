package pworm

type FakeExecuter struct {
}

func (f *FakeExecuter) Execute(command string) ([]byte, error) {
	return nil, nil
}
