package condition

import (
	"fmt"
	"strings"
)

type Ends struct {
	InlineString
}

func (e *Ends) Describe(value string) string {
	return fmt.Sprintf("(%s) must be ends with %q", value, e.value)
}

func (e *Ends) Check(value string) bool {
	return strings.HasSuffix(value, e.value)
}
