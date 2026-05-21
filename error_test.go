package pworm

import (
	"errors"
	"testing"
)

func TestErrorWrap(t *testing.T) {

	errNotFound := errors.New("not found")

	var err = NewError(errNotFound, "object")

	if !errors.Is(err, errNotFound) {
		t.Errorf("Expected error to be %v, got %v", errNotFound, err)
	}

	if err.Error() != "not found: object" {
		t.Errorf("Expected error message to be %s, got %s", "not found: object", err.Error())
	}
}
