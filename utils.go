package main

import (
	"fmt"

	"github.com/gotify/plugin-api"
)

func makeMarkdownMessage(title, message string, clickURL *string) plugin.Message {
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
		Title:   title,
		Message: message,
		Extras:  extras,
	}
}

func makeClickURL(messageID int, host, org, name, appendix string) *string {
	clickURLTeml := "%s/org/%s/servers/%s/messages/%d%s" // host, org, server name, message ID, appendix
	s := fmt.Sprintf(clickURLTeml, host, org, name, messageID, appendix)
	return &s
}
