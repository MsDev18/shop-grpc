package product

import "regexp"

var SLUG_REGEX = regexp.MustCompile(`^[a-z0-9]+(-[a-z0-9]+)*$`)

const (
	NAME_MIN_LENGTH = 3
	NAME_MAX_LENGTH = 255
	SLUG_MIN_LENGTH = 3
	SLUG_MAX_LENGTH = 255
)

type Validator struct {

}

func New () Validator {
	return Validator{
		
	}
}