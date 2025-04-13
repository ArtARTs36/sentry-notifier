package mattermostwh

import (
	"fmt"
	"net/url"

	"gopkg.in/yaml.v3"
)

type urlWrapper struct {
	url.URL
}

func (w *urlWrapper) UnmarshalYAML(n *yaml.Node) error {
	if n.Kind != yaml.ScalarNode {
		return fmt.Errorf("must be a string, got %q", n.Kind)
	}

	value, err := url.Parse(n.Value)
	if err != nil {
		return fmt.Errorf("invalid url: %w", err)
	}

	w.URL = *value

	return err
}
