package condition

type String struct {
	Equals   *Equals   `yaml:"equals,omitempty" json:"equals,omitempty"`
	Starts   *Starts   `yaml:"starts,omitempty" json:"starts,omitempty"`
	Ends     *Ends     `yaml:"ends,omitempty" json:"ends,omitempty"`
	Contains *Contains `yaml:"contains,omitempty" json:"contains,omitempty"`
}

func (s *String) Check(value string) (bool, string) {
	if s.Equals != nil && !s.Equals.Check(value) {
		return false, s.Equals.Describe(value)
	}

	if s.Starts != nil && !s.Starts.Check(value) {
		return false, s.Starts.Describe(value)
	}

	if s.Ends != nil && !s.Ends.Check(value) {
		return false, s.Ends.Describe(value)
	}

	if s.Contains != nil && !s.Contains.Check(value) {
		return false, s.Contains.Describe(value)
	}

	return true, ""
}
