package richerror

import "fmt"

type RichError struct {
	op      string
	message string
	kind    Kind
	err     error
	meta    map[string]any
}

func New() *RichError {
	return &RichError{}
}

func (r *RichError) Error() string {
	errMsg := fmt.Sprintf("%s -> %s", r.op, r.GetMessage())
	return errMsg
}

func (r *RichError) Unwrap() error {
	return r.err
}
