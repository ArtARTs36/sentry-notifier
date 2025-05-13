package condition

import (
	"fmt"
	"strings"
)

type Contains struct {
	InlineString
}

func (e *Contains) Describe(value string) string {
	return fmt.Sprintf("(%s) must be contains %q", value, e.value)
}

func (e *Contains) Check(value string) bool {
	return strings.HasSuffix(value, e.value)
}
