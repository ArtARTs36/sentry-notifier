package sentry

import "encoding/json"

type IssuePayload struct {
	Action string `json:"action"`
	Data   struct {
		Issue struct {
			Count   string `json:"count"`
			ID      string `json:"id"`
			Level   string `json:"level"`
			ShortID string `json:"shortId"`
			Status  string `json:"status"`
			Type    string `json:"type"`
			Title   string `json:"title"`
			Project struct {
				ID       string `json:"id"`
				Name     string `json:"name"`
				Platform string `json:"platform"`
				Slug     string `json:"slug"`
			} `json:"project"`
			LastSeen  Time `json:"lastSeen"`
			FirstSeen Time `json:"firstSeen"`
		} `json:"issue"`
	} `json:"data"`
}

func createIssuePayload(data []byte) (*IssuePayload, error) {
	pl := new(IssuePayload)
	err := json.Unmarshal(data, &pl)
	if err != nil {
		return nil, err
	}

	return pl, nil
}

func (p *IssuePayload) GetHookResource() HookResource {
	return HookResourceIssue
}

func (p *IssuePayload) GetData() interface{} {
	return p.Data
}
