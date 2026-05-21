package pworm

import "github.com/anCreny/pw-orm/errors"

type result struct {
	Err *errors.Error `json:"Error"`
	Out []byte        `json:"Output"`
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
	return r.Out
}
