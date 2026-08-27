package richerror

func (r *RichError) SetOp(op string) *RichError {
	r.op = op
	return r
}

func (r *RichError) SetMsg(message string) *RichError {
	r.message = message
	return r
}

func (r *RichError) SetKind(kind Kind) *RichError {
	r.kind = kind
	return r
}

func (r *RichError) SetErr(err error) *RichError {
	r.err = err
	return r
}

func (r *RichError) SetMeta(meta map[string]any) *RichError {
	r.meta = meta
	return r
}
