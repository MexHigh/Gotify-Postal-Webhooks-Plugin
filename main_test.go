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

var messageBouncedEvent = []byte(`{
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

var messageLinkClickedEvent = []byte(`{
	"event": "MessageLinkClicked",
	"timestamp": 0.0,
	"uuid": "irrelevant",
	"payload": {
		"url":"https://atech.media",
		"token":"VJzsFA0S",
		"ip_address":"185.22.208.2",
		"user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36",
		"message":{
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

var messageLoadedEvent = []byte(`{
	"event": "MessageLoaded",
	"timestamp": 0.0,
	"uuid": "irrelevant",
	"payload": {
		"ip_address":"185.22.208.2",
		"user_agent":"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36",
		"message":{
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

var domainDNSErrorEvent = []byte(`{
	"event": "DomainDNSError",
	"timestamp": 0.0,
	"uuid": "irrelevant",
	"payload": {
		"domain":"example.com",
		"uuid":"820b47a4-4dfd-42e4-ae6a-1e5bed5a33fd",
		"dns_checked_at":1477945711.5502,
		"spf_status":"OK",
		"spf_error":null,
		"dkim_status":"Invalid",
		"dkim_error":"The DKIM record at example.com does not match the record we have provided. Please check it has been copied correctly.",
		"mx_status":"Missing",
		"mx_error":null,
		"return_path_status":"OK",
		"return_path_error":null,
		"server":{
			"uuid":"54529725-8807-4069-ab29-a3746c1bbd98",
			"name":"AwesomeApp Mail Server",
			"permalink":"awesomeapp",
			"organization":"atech"
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

func TestProcessWebhookMessageBouncedTitle(t *testing.T) {
	p := &Plugin{
		config: &PluginConfig{},
	}
	result := p.processWebhookBytes(messageBouncedEvent, nil)
	mdMsg := makeMarkdownMessage(result.Title, result.Message, result.clickURL)

	if mdMsg.Title != EmojiExclamMark+" Bounce message received" {
		t.Fatal("Message title does not match, got: ", mdMsg.Title)
	}
}

func TestProcessWebhookMessageLinkClickedTitle(t *testing.T) {
	p := &Plugin{
		config: &PluginConfig{},
	}
	result := p.processWebhookBytes(messageLinkClickedEvent, nil)
	mdMsg := makeMarkdownMessage(result.Title, result.Message, result.clickURL)

	if mdMsg.Title != EmojiEyes+" Link in message was clicked" {
		t.Fatal("Message title does not match, got: ", mdMsg.Title)
	}
}

func TestProcessWebhookMessageLoadedTitle(t *testing.T) {
	p := &Plugin{
		config: &PluginConfig{},
	}
	result := p.processWebhookBytes(messageLoadedEvent, nil)
	mdMsg := makeMarkdownMessage(result.Title, result.Message, result.clickURL)

	if mdMsg.Title != EmojiEyes+" Message was opened" {
		t.Fatal("Message title does not match, got: ", mdMsg.Title)
	}
}

func TestProcessWebhookDomainDNSErrorTitle(t *testing.T) {
	p := &Plugin{
		config: &PluginConfig{},
	}
	result := p.processWebhookBytes(domainDNSErrorEvent, nil)
	mdMsg := makeMarkdownMessage(result.Title, result.Message, result.clickURL)

	if mdMsg.Title != EmojiExclamMark+" DNS setup check failed" {
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
