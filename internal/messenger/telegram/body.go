package telegram

import (
	"strings"

	"github.com/artarts36/sentry-notifier/internal/messenger/contracts"
)

var reservedChars = []string{
	"\\", "_", "*", "[", "]",
	"(", ")", "~", "!", ">",
	"<", "&", "#", "+", "-",
	"=", "|", "{", "}", ".",
}

func (t *Telegram) prepareBody(message contracts.Message) string {
	body := message.Body

	for _, char := range reservedChars {
		body = strings.ReplaceAll(body, char, "\\"+char)
	}

	return body
}
