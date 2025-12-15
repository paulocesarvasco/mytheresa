package catalogapi

import "regexp"

var codeValidator = regexp.MustCompile(`^PROD\d{3}$`)

func isValidProductCode(code string) bool {
	return codeValidator.MatchString(code)
}
