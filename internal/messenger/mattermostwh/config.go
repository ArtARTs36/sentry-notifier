package mattermostwh

import (
	"github.com/artarts36/specw"
)

type Config struct {
	URL specw.Env[specw.URL] `yaml:"url" json:"url"`
}
