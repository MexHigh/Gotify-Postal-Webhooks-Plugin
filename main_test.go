package main

import (
	"strings"
	"testing"

	"github.com/gotify/plugin-api"
)

var messageSentEvent = []byte(`{
	"event": "MessageSent",
	"timestamp": 0.0,
	"uuid": "irrelevant",
	"payload": {
		"status":"Sent",
		"details":"Message sent by SMTP to aspmx.l.google.com (2a00:1450:400c:c0b::1b) (from 2a00:67a0:a:15::2)",
		"output":"250 2.0.0 OK 1477944899 ly2si31746747wjb.95 - gsmtp",
		"time":0.22,
		"sent_with_ssl":true,
		"timestamp":1477945177.12994,
		"message": {
			"id":12345,
			"token":"abcdef123",
			"direction":"outgoing",
			"message_id":"5817a64332f44_4ec93ff59e79d154565eb@app34.mail",
			"to":"test@example.com",
			"from":"sales@awesomeapp.com",
			"subject":"Welcome to AwesomeApp",
			"timestamp":1477945177.12994,
			"spam_status":"NotSpam",
			"tag":"welcome"
		}
	}
}`)

var messageBounced = []byte(`{
	"event": "MessageBounced",
	"timestamp": 0.0,
	"uuid": "irrelevant",
	"payload": {
		"original_message":{
			"id":12345,
			"token":"abcdef123",
			"direction":"outgoing",
			"message_id":"5817a64332f44_4ec93ff59e79d154565eb@app34.mail",
			"to":"test@example.com",
			"from":"sales@awesomeapp.com",
			"subject":"Welcome to AwesomeApp",
			"timestamp":1477945177.12994,
			"spam_status":"NotSpam",
			"tag":"welcome"
		},
		"bounce":{
			"id":12347,
			"token":"abcdef124",
			"direction":"incoming",
			"message_id":"5817a64332f44_4ec93ff59e79d154565eb@someserver.com",
			"to":"abcde@psrp.postal.yourdomain.com",
			"from":"postmaster@someserver.com",
			"subject":"Delivery Error",
			"timestamp":1477945179.12994,
			"spam_status":"NotSpam",
			"tag":null
		}
	}
}`)

/*type MockedMessageHandler struct {
	Chan chan plugin.Message
}

func (m MockedMessageHandler) SendMessage(msg plugin.Message) error {
	fmt.Println(msg.Title)
	m.Chan <- msg
	return nil
}*/

func TestProcessWebhookWithClickURL(t *testing.T) {
	p := &Plugin{
		config: &PluginConfig{},
	}
	msInfo := &PostalMailserverInfo{
		Host:         "https://testing.example.com",
		Organization: "testing-org",
		Name:         "testing-server",
	}

	result := p.processWebhookBytes(messageSentEvent, msInfo)
	mdMsg := makeMarkdownMessage(result.Title, result.Message, result.clickURL)

	if !hasClickURL(mdMsg) {
		t.Fatal("Message doesn't have a click URL")
	}
	expectedPrefix := "https://testing.example.com/org/testing-org/servers/testing-server/messages/"
	if cu := getClickURL(mdMsg); !strings.HasPrefix(cu, expectedPrefix) {
		t.Fatal("Message has wrong clickURL prefix, expected prefix: "+expectedPrefix, ", got: "+cu)
	}
}

func TestProcessWebhookWithoutClickURL(t *testing.T) {
	p := &Plugin{
		config: &PluginConfig{},
	}

	result := p.processWebhookBytes(messageSentEvent, nil)
	mdMsg := makeMarkdownMessage(result.Title, result.Message, result.clickURL)

	if hasClickURL(mdMsg) {
		t.Fatal("Message has clickURL which it shouldn't have")
	}
}

func TestProcessWebhookMessageSentTitle(t *testing.T) {
	p := &Plugin{
		config: &PluginConfig{},
	}
	result := p.processWebhookBytes(messageSentEvent, nil)
	mdMsg := makeMarkdownMessage(result.Title, result.Message, result.clickURL)

	if mdMsg.Title != EmojiCheckMark+" Message delivered successfully" {
		t.Fatal("Message title does not match, got: ", mdMsg.Title)
	}
}

// Utilitiy test functions

func hasClickURL(msg plugin.Message) bool {
	return getClickURL(msg) != ""
}

func getClickURL(msg plugin.Message) string {
	if notif, ok := msg.Extras["client::notification"]; ok {
		if notifMap, ok := notif.(map[string]interface{}); ok {
			if click, ok := notifMap["click"]; ok {
				if clickMap, ok := click.(map[string]string); ok {
					if url, ok := clickMap["url"]; ok {
						return url
					}
				}
			}
		}
	}
	return ""
}
