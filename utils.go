package main

import (
	"fmt"

	"github.com/gotify/plugin-api"
)

const (
	messageTemplate          = "Sender: %s\n\n%s"
	messageTemplateCodeBlock = "Sender: %s\n\n```\n%s\n```"
)

func makeMarkdownMessage(title, message, remoteIP string, withinCodeBlock bool) plugin.Message {
	tmpl := messageTemplate
	if withinCodeBlock {
		tmpl = messageTemplateCodeBlock
	}

	return plugin.Message{
		Title: title,
		Message: fmt.Sprintf(tmpl,
			remoteIP,
			message,
		),
		Extras: map[string]interface{}{
			"client::display": map[string]interface{}{
				"contentType": "text/markdown",
			},
		},
	}
}
