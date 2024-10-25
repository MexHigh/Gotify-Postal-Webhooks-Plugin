package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

// TODO add emojis

func (p *Plugin) handleMessageStatusEvent(payload json.RawMessage, eventType WebhookMessageEvent) (*GotifyMessage, error) {
	var msg MessageStatusEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	message.clickURL = p.config.makeClickURL(msg.Message.ID, "")

	// message status events can have several eventTypes, so we need to switch here again
	switch eventType {
	case WebhookMessageEventMessageSent:
		message.Title = fmt.Sprintf("Message for %s delivered successfully", msg.Message.To)
	case WebhookMessageEventMessageDelayed:
		message.Title = fmt.Sprintf("Message for %s was delayed", msg.Message.To)
	case WebhookMessageEventMessageDeliveryFailed:
		message.Title = fmt.Sprintf("Message delivery for %s failed", msg.Message.To)
	case WebhookMessageEventMessageHeld:
		message.Title = fmt.Sprintf("Message delivery to %s was held", msg.Message.To)
	default:
		return nil, errors.New("unknown event name '" + string(eventType) + "' occured in message status event handler")
	}

	message.Message += fmt.Sprintf("_%s &rarr; %s: \"%s\"_\n\n", msg.Message.From, msg.Message.To, msg.Message.Subject)
	message.Message += fmt.Sprintf("%s\n\n---\n\n", msg.Details)
	message.Message += fmt.Sprintf("**Delivery time:** %f seconds\n\n", msg.Time)
	message.Message += fmt.Sprintf("**Sent with SSL/TLS:** %t\n\n", msg.SentWithSSL)
	message.Message += fmt.Sprintf("**Output:**\n\n```%s```", msg.Output)

	return message, nil
}

func (p *Plugin) handleMessageLoadedEvent(payload json.RawMessage) (*GotifyMessage, error) {
	var msg MessageLoadedEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	message.clickURL = p.config.makeClickURL(msg.Message.ID, "/activity")

	// TODO
	message.Title = "not implemented"

	return message, nil
}

func (p *Plugin) handleMessageBounceEvent(payload json.RawMessage) (*GotifyMessage, error) {
	var msg MessageBounceEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	message.clickURL = p.config.makeClickURL(msg.OriginalMessage.ID, "")

	// TODO
	message.Title = "not implemented"

	return message, nil
}

func (p *Plugin) handleMessageClickEvent(payload json.RawMessage) (*GotifyMessage, error) {
	var msg MessageClickEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	message.clickURL = p.config.makeClickURL(msg.Message.ID, "/activity")

	// TODO
	message.Title = "not implemented"

	return message, nil
}

func (p *Plugin) handleDNSErrorEvent(payload json.RawMessage) (*GotifyMessage, error) {
	var msg DNSErrorEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	//message.clickURL = p.config.makeClickURL(msg.Message.ID)

	// TODO
	message.Title = "not implemented"

	return message, nil
}
