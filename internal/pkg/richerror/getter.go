package richerror


func (r *RichError) GetMessage() string {
	if r.message != "" {
		return r.message
	}
	if r.err == nil {
		return ""
	}
	if rErr, ok := r.err.(*RichError); ok {
		return rErr.GetMessage()
	}
	return r.err.Error()
}

func (r *RichError) GetKind() Kind {
	if r.kind != 0 {
		return r.kind
	}
	if r.err == nil {
		return 0
	}
	if rErr, ok := r.err.(*RichError) ; ok {
		return rErr.GetKind()
	}
	return 0
}


func (r *RichError) GetMeta () map[string]any {
	return r.meta
}