package condition

import "fmt"

type Equals struct {
	InlineString
}

func (e *Equals) Describe(value string) string {
	return fmt.Sprintf("(%s) must be equals %q", value, e.value)
}

func (e *Equals) Check(value string) bool {
	return value == e.value
}
