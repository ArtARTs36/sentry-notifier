package condition

import (
	"fmt"
	"strings"
)

type Starts struct {
	InlineString
}

func (s *Starts) Describe(value string) string {
	return fmt.Sprintf("(%s) must be starts with %q", value, s.value)
}

func (s *Starts) Check(value string) bool {
	return strings.HasPrefix(value, s.value)
}
