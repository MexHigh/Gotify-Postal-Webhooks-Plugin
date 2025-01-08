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
		ModulePath:  "git.leon.wtf/leon/gotify-postal-webhooks-plugin",
		Author:      "Leon Schmidt <mail@leon-schmidt.dev>",
		Website:     "https://leon-schmidt.dev",
	}
}

type GotifyMessage struct {
	Title    string
	Message  string
	clickURL *string
}

type PostalMailserverInfo struct {
	Host         string
	Organization string
	Name         string
}

type PluginConfig struct {
	VerboseOutput bool
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
	return &PluginConfig{false}
}

// ValidateAndSetConfig implements plugin.Configurer
func (p *Plugin) ValidateAndSetConfig(c interface{}) error {
	config := c.(*PluginConfig)
	p.config = config
	return nil
}

const helpMessageTemplate = "Use this **webhook URL**: %s\n\n" +
	"You can also set the Postal host, organization and server name as parameters (e.g. `?host=postal.example.com&org=some-org&name=main`). " +
	"Once done, Gotify messages can be clicked to open the corresponding dashboard in Postal."

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
				nil,
			))
			return
		}

		if p.config.VerboseOutput {
			fmt.Println("Incoming Postal webhook:")
			fmt.Println(string(bytes))
		}

		// get optional params
		var msInfo *PostalMailserverInfo
		host, hostExists := c.Params.Get("host")
		org, orgExists := c.Params.Get("org")
		name, nameExists := c.Params.Get("name")
		if hostExists && orgExists && nameExists {
			msInfo = &PostalMailserverInfo{
				Host:         host,
				Organization: org,
				Name:         name,
			}
		}

		// this function does not return error since errors are handled within
		// the function and returned "pre-serialized" as GotifyMessages
		notification := p.processWebhookBytes(bytes, msInfo)

		// send message
		p.msgHandler.SendMessage(makeMarkdownMessage(
			notification.Title,
			notification.Message,
			notification.clickURL, // may be nil
		))
	}

	mux.POST("/"+routeName, webhookHandler)
}

func (p *Plugin) processWebhookBytes(bytes []byte, msInfo *PostalMailserverInfo) *GotifyMessage {
	// unmarshal body to generic WebhookMessage
	var message WebhookMessage
	if err := json.Unmarshal(bytes, &message); err != nil {
		return &GotifyMessage{
			"Error unmarshalling Postal message",
			err.Error(),
			nil,
		}
	}

	// switch message event type for parsing PayloadRaw
	switch message.Event {
	// all message status events
	case WebhookMessageEventMessageSent, WebhookMessageEventMessageDelayed, WebhookMessageEventMessageDeliveryFailed, WebhookMessageEventMessageHeld:
		notification, err := p.handleMessageStatusEvent(message.PayloadRaw, message.Event, msInfo)
		if err != nil {
			return &GotifyMessage{
				"Error handling message status event",
				err.Error(),
				nil,
			}
		}
		return notification

	// bounce events
	case WebhookMessageEventMessageBounced:
		notification, err := p.handleMessageBounceEvent(message.PayloadRaw, msInfo)
		if err != nil {
			return &GotifyMessage{
				"Error handling message status event",
				err.Error(),
				nil,
			}
		}
		return notification

	// linktracking link clicked
	case WebhookMessageEventMessageLinkClicked:
		notification, err := p.handleMessageClickEvent(message.PayloadRaw, msInfo)
		if err != nil {
			return &GotifyMessage{
				"Error handling message status event",
				err.Error(),
				nil,
			}
		}
		return notification

	// message loaded events
	case WebhookMessageEventMessageLoaded:
		notification, err := p.handleMessageLoadedEvent(message.PayloadRaw, msInfo)
		if err != nil {
			return &GotifyMessage{
				"Error handling message status event",
				err.Error(),
				nil,
			}
		}
		return notification

	// DNS error
	case WebhookMessageEventDomainDNSError:
		notification, err := p.handleDNSErrorEvent(message.PayloadRaw)
		if err != nil {
			return &GotifyMessage{
				"Error handling message status event",
				err.Error(),
				nil,
			}
		}
		return notification

	default:
		return &GotifyMessage{
			"Read unknown event name in Postal massage",
			fmt.Sprintf("Event name was '%s'", string(message.Event)),
			nil,
		}
	}
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
