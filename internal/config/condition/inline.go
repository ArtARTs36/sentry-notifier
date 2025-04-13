package condition

import (
	"fmt"

	"gopkg.in/yaml.v3"
)

type InlineString struct {
	value string
}

func (s *InlineString) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("must be a string, got %q", n.Kind)
	}

	s.value = n.Value

	return nil
}

func (s *InlineString) UnmarshalJSON(data []byte) error {

}
