package slugvalidator

import (
	"strings"

	"github.com/gosimple/slug"
)

func IsSlug(slugString string) bool {
	slugString = strings.ToLower(slugString)
	return slug.IsSlug(slugString)
}
