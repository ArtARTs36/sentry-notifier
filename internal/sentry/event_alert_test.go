package sentry

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestEventAlert_extract(t *testing.T) {
	al := &EventAlert{}

	al.Data.Event.URL = "https://sentry.io/api/0/projects/test-org/front-end/events/e4874d664c3540c1a32eab185f12c5ab/"

	al.extract()

	assert.Equal(t, "test-org", al.Data.Extracted.OrganizationName)
	assert.Equal(t, "front-end", al.Data.Extracted.ProjectName)
}
