package category

import "regexp"

type Validator struct {
}

func New() Validator {
	return Validator{}
}

var SLUG_REGEX = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

const (
	TITLE_MIN_LENGTH = 3
	TITLE_MAX_LENGTH = 50
	SLUG_MIN_LENGTH  = 3
	SLUG_MAX_LENGTH = 255
)