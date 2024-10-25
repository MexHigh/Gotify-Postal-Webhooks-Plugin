package main

import (
	"encoding/json"
	"time"
)

type WebhookMessageEvent string

const (
	// send events
	WebhookMessageEventMessageSent           WebhookMessageEvent = "MessageSent"
	WebhookMessageEventMessageDelayed        WebhookMessageEvent = "MessageDelayed"
	WebhookMessageEventMessageDeliveryFailed WebhookMessageEvent = "MessageDeliveryFailed"
	WebhookMessageEventMessageHeld           WebhookMessageEvent = "MessageHeld"

	// other events
	WebhookMessageEventMessageLoaded      WebhookMessageEvent = "MessageLoaded"
	WebhookMessageEventMessageBounced     WebhookMessageEvent = "MessageBounced"
	WebhookMessageEventMessageLinkClicked WebhookMessageEvent = "MessageLinkClicked"
	WebhookMessageEventDomainDNSError     WebhookMessageEvent = "DomainDNSError"
)

type WebhookMessage struct {
	Event      WebhookMessageEvent `json:"event"`
	Timestamp  time.Time           `json:"timestamp"`
	UUID       string              `json:"uuid"`
	PayloadRaw json.RawMessage     `json:"payload"`
}

type MessageStatusEvent struct {
	Status      string  `json:"status"`
	Details     string  `json:"details"`
	Output      string  `json:"output"`
	Time        float64 `json:"time"`
	SentWithSSL bool    `json:"sent_with_ssl"`
	Timestamp   float64 `json:"timestamp"`
	Message     Message `json:"message"`
}

type MessageBounceEvent struct {
	OriginalMessage Message `json:"original_message"`
	Bounce          Message `json:"bounce"`
}

type MessageClickEvent struct {
	URL       string  `json:"url"`
	Token     string  `json:"token"`
	IPAddress string  `json:"ip_address"`
	UserAgent string  `json:"user_agent"`
	Message   Message `json:"message"`
}

type MessageLoadedEvent struct {
	IPAddress string  `json:"ip_address"`
	UserAgent string  `json:"user_agent"`
	Message   Message `json:"message"`
}

type DNSErrorEvent struct {
	Domain           string  `json:"domain"`
	UUID             string  `json:"uuid"`
	DNSCheckedAt     float64 `json:"dns_checked_at"`
	SPFStatus        string  `json:"spf_status"`
	SPFError         string  `json:"spf_error"`
	DKIMStatus       string  `json:"dkim_status"`
	DKIMError        string  `json:"dkim_error"`
	MXStatus         string  `json:"mx_status"`
	MXError          string  `json:"mx_error"`
	ReturnPathStatus string  `json:"return_path_status"`
	ReturnPathError  string  `json:"return_path_error"`
	Server           Server  `json:"server"`
}

type Message struct {
	ID         int     `json:"id"`
	Token      string  `json:"token"`
	Direction  string  `json:"direction"`
	MessageID  string  `json:"message_id"`
	To         string  `json:"to"`
	From       string  `json:"from"`
	Subject    string  `json:"subject"`
	Timestamp  float64 `json:"timestamp"`
	SpamStatus string  `json:"spam_status"`
	Tag        *string `json:"tag"`
}

type Server struct {
	UUID         string `json:"uuid"`
	Name         string `json:"name"`
	Permalink    string `json:"permalink"`
	Organization string `json:"organization"`
}
