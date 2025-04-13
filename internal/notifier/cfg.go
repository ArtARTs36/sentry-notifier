package notifier

import (
	"fmt"

	"github.com/artarts36/sentry-notifier/internal/config/condition"
	"github.com/artarts36/sentry-notifier/internal/sentry"
)

type Strategy string

const (
	StrategyDefault     Strategy = StrategyImmediately
	StrategyImmediately Strategy = "immediately"
	StrategyAsync       Strategy = "async"
	StrategyNull        Strategy = "null"
)

type Config struct {
	Strategy Strategy                           `yaml:"strategy" json:"strategy"`
	On       map[sentry.HookResource][]Template `yaml:"on" json:"on"`
}

type Template struct {
	When    *TemplateWhen `yaml:"when" json:"when"`
	Message string        `yaml:"message" json:"message"`
	To      string        `yaml:"to" json:"to"`

	MessageTemplateID string `yaml:"-" json:"-"`
}

type TemplateWhen struct {
	ProjectName condition.String `yaml:"project_name" json:"project_name"`
}

func (t *TemplateWhen) Check(pl sentry.Payload) (bool, string) {
	if t == nil {
		return true, ""
	}

	if state, reason := t.ProjectName.Check(pl.GetProjectSlug()); !state {
		return false, fmt.Sprintf("project_name %s", reason)
	}

	return true, ""
}
