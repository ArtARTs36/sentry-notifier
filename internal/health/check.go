package health

type CheckResponse struct {
	Status bool `json:"status"`

	Checks []Check `json:"checks"`
}

type Check struct {
	ID     string `json:"id"`
	Status bool   `json:"status"`
	Reason string `json:"reason,omitempty"`
}
