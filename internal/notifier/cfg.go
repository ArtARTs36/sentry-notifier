package notifier

import "github.com/artarts36/sentry-notifier/internal/sentry"

type Strategy string

const (
	StrategyDefault     Strategy = StrategyImmediately
	StrategyImmediately Strategy = "immediately"
	StrategyAsync       Strategy = "async"
)

type Config struct {
	Strategy Strategy                           `yaml:"strategy"`
	On       map[sentry.HookResource][]Template `yaml:"on"`
}
