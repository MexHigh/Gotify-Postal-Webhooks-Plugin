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

type GotifyMessage struct {
	Title    string
	Message  string
	clickURL *string // optional
}

type PluginConfig struct {
	Host           *string
	Organization   *string
	MailserverName *string
}

func (pc *PluginConfig) makeClickURL(messageID int, appendix string) *string {
	if pc.Host != nil && pc.Organization != nil && pc.MailserverName != nil {
		clickURLTeml := "%s/org/%s/servers/%s/messages/%d%s" // host, org, server name, message ID, appendix
		url := fmt.Sprintf(clickURLTeml, pc.Host, pc.Organization, pc.MailserverName, messageID, appendix)
		return &url
	} else {
		return nil
	}
}

// Plugin is plugin instance
type Plugin struct {
	userCtx    plugin.UserContext
	msgHandler plugin.MessageHandler
	basePath   string
	config     *PluginConfig
}

// Enable implements plugin.Plugin
func (p *Plugin) Enable() error {
	return nil
}

// Disable implements plugin.Plugin
func (p *Plugin) Disable() error {
	return nil
}

// DefaultConfig implements plugin.Configurer
func (p *Plugin) DefaultConfig() interface{} {
	return &PluginConfig{nil, nil, nil}
}

// ValidateAndSetConfig implements plugin.Configurer
func (p *Plugin) ValidateAndSetConfig(c interface{}) error {
	config := c.(*PluginConfig)
	p.config = config
	return nil
}

const helpMessageTemplate = "Use this **webhook URL**: %s\n\n" +
	"You can also configure the Postal instance and server here. Once done, Gotify messages can be clicked to open the corresponding dashboard in Postal."

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
				nil,
			))
			return
		}

		// unmarshal body to generic WebhookMessage
		var message WebhookMessage
		if err := json.Unmarshal(bytes, &message); err != nil {
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Error unmarshalling Postal message",
				err.Error(),
				c.ClientIP(),
				nil,
			))
			return
		}

		notification := &GotifyMessage{}

		switch message.Event {
		// all message status events
		case WebhookMessageEventMessageSent, WebhookMessageEventMessageDelayed, WebhookMessageEventMessageDeliveryFailed, WebhookMessageEventMessageHeld:
			notification, err = p.handleMessageStatusEvent(message.PayloadRaw, message.Event)
			if err != nil {
				p.msgHandler.SendMessage(makeMarkdownMessage(
					"Error handling message status event",
					err.Error(),
					c.ClientIP(),
					nil,
				))
				return
			}

		// message loaded events
		case WebhookMessageEventMessageLoaded:
			notification, err = p.handleMessageLoadedEvent(message.PayloadRaw)
			if err != nil {
				p.msgHandler.SendMessage(makeMarkdownMessage(
					"Error handling message status event",
					err.Error(),
					c.ClientIP(),
					nil,
				))
				return
			}

		// bounce events
		case WebhookMessageEventMessageBounced:
			notification, err = p.handleMessageBounceEvent(message.PayloadRaw)
			if err != nil {
				p.msgHandler.SendMessage(makeMarkdownMessage(
					"Error handling message status event",
					err.Error(),
					c.ClientIP(),
					nil,
				))
				return
			}

		// linktracking link clicked
		case WebhookMessageEventMessageLinkClicked:
			notification, err = p.handleMessageClickEvent(message.PayloadRaw)
			if err != nil {
				p.msgHandler.SendMessage(makeMarkdownMessage(
					"Error handling message status event",
					err.Error(),
					c.ClientIP(),
					nil,
				))
				return
			}

		// DNS error
		case WebhookMessageEventDomainDNSError:
			notification, err = p.handleDNSErrorEvent(message.PayloadRaw)
			if err != nil {
				p.msgHandler.SendMessage(makeMarkdownMessage(
					"Error handling message status event",
					err.Error(),
					c.ClientIP(),
					nil,
				))
				return
			}

		default:
			p.msgHandler.SendMessage(makeMarkdownMessage(
				"Read unknown event name in Postal massage",
				fmt.Sprintf("Event name was '%s'", string(message.Event)),
				c.ClientIP(),
				nil,
			))
			return
		}

		// send final message
		p.msgHandler.SendMessage(makeMarkdownMessage(
			notification.Title,
			notification.Message,
			c.ClientIP(),
			notification.clickURL, // may be nil
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
