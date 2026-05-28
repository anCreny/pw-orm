package pworm

import (
	"encoding/json"

	"github.com/anCreny/pw-orm/errors"
)

type result struct {
	Err *errors.Error    `json:"Error"`
	Out *json.RawMessage `json:"Output"`
}

func (r *result) FullError() *errors.Error {
	return r.Err
}

func (r *result) Error() error {
	if r.Err != nil {
		categoryInfo := r.Err.CategoryInfo
		err := errors.GetCategoryEnumItem(categoryInfo.Category)

		message := r.Err.Exception.Message

		return NewError(err, message)
	}

	return nil
}

func (r *result) Output() []byte {
	if r.Out == nil {
		return nil
	}
	return *r.Out
}
