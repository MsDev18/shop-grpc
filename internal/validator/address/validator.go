package address

import "regexp"

var POSTAL_CODE_REGEX = regexp.MustCompile(`^\d{10}$`)
var PROVINCE_ID_REGEX = regexp.MustCompile(`^[0-9]+$`)

type Validator struct {
}

func New() Validator {
	return Validator{}
}
