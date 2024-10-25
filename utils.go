package main

import (
	"fmt"

	"github.com/gotify/plugin-api"
)

const (
	messageTemplate = "Sender: %s\n\n%s"
)

func makeMarkdownMessage(title, message, remoteIP string, clickURL *string) plugin.Message {
	tmpl := messageTemplate

	extras := map[string]interface{}{}
	extras["client::display"] = map[string]interface{}{
		"contentType": "text/markdown",
	}
	if clickURL != nil {
		extras["client::notification"] = map[string]interface{}{
			"click": map[string]string{ // map[string]interface{} ?
				"url": *clickURL,
			},
		}
	}

	return plugin.Message{
		Title: title,
		Message: fmt.Sprintf(tmpl,
			remoteIP,
			message,
		),
		Extras: extras,
	}
}
