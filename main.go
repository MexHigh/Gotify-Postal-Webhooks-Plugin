package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/gotify/plugin-api"
)

const routeName = "postal"

// GetGotifyPluginInfo returns gotify plugin info
func GetGotifyPluginInfo() plugin.Info {
	return plugin.Info{
		Name:        "Postal Webhooks",
		Description: "Plugin that enables Gotify to receive webhooks from Postal",
		ModulePath:  "git.leon.wtf/leon/gotify-postal-webhook-plugin",
		Author:      "Leon Schmidt <mail@leon-schmidt.dev>",
		Website:     "https://leon-schmidt.dev",
	}
}

// Plugin is plugin instance
type Plugin struct {
	userCtx    plugin.UserContext
	msgHandler plugin.MessageHandler
	basePath   string
}

// Enable implements plugin.Plugin
func (p *Plugin) Enable() error {
	return nil
}

// Disable implements plugin.Plugin
func (p *Plugin) Disable() error {
	return nil
}

const helpMessageTemplate = "Use this **webhook URL**: %s"

// GetDisplay implements plugin.Displayer
func (p *Plugin) GetDisplay(location *url.URL) string {
	baseHost := ""
	if location != nil {
		baseHost = fmt.Sprintf("%s://%s", location.Scheme, location.Host)
	}
	webhookURL := baseHost + p.basePath + routeName
	return fmt.Sprintf(helpMessageTemplate, webhookURL)
}

// SetMessageHandler implements plugin.Messenger
func (p *Plugin) SetMessageHandler(h plugin.MessageHandler) {
	// invoced during initialization
	p.msgHandler = h
}

// RegisterWebhook implements plugin.Webhooker
func (p *Plugin) RegisterWebhook(basePath string, mux *gin.RouterGroup) {
	p.basePath = basePath

	webhookHandler := func(c *gin.Context) {
		// read body
		bytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Error reading request body",
				err.Error(),
				c.ClientIP(),
				false,
			))
			return
		}

		// unmarshal body
		var message WebhookMessage
		if err := json.Unmarshal(bytes, &message); err != nil {
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Error unmarshalling Postal message",
				err.Error(),
				c.ClientIP(),
				false,
			))
			return
		}

		notificationTitle := ""
		notificationMessage := ""

		switch message.Event {
		case WebhookMessageEventMessageSent,
			WebhookMessageEventMessageDelayed,
			WebhookMessageEventMessageDeliveryFailed,
			WebhookMessageEventMessageHeld: // All message events
			// TODO
		case WebhookMessageEventMessageLoaded:
			// TODO
		case WebhookMessageEventMessageBounced:
			// TODO
		case WebhookMessageEventMessageLinkClicked:
			// TODO
		case WebhookMessageEventDomainDNSError:
			// TODO
		default:
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Read unknown event name in Postal massage",
				fmt.Sprintf("Event name was '%s'", string(message.Event)),
				c.ClientIP(),
				false,
			))
			return
		}

		// send final message
		p.msgHandler.SendMessage(makeMarkdownMessage(
			notificationTitle,
			notificationMessage,
			c.ClientIP(),
			false,
		))
	}

	mux.POST("/"+routeName, webhookHandler)
}

// NewGotifyPluginInstance creates a plugin instance for a user context.
func NewGotifyPluginInstance(ctx plugin.UserContext) plugin.Plugin {
	return &Plugin{
		userCtx: ctx,
	}
}

func main() {
	panic("this should be built as go plugin")
}
