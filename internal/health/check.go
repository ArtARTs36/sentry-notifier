package health

type Check struct {
	Status bool `json:"status"`

	Channels map[string]map[string][]CheckChannel `json:"channels"`
}

type CheckChannel struct {
	Status bool   `json:"status"`
	Reason string `json:"reason,omitempty"`
}
