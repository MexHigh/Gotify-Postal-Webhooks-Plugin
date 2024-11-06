package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

// TODO add emojis

func (p *Plugin) handleMessageStatusEvent(payload json.RawMessage, eventType WebhookMessageEvent, msInfo *PostalMailserverInfo) (*GotifyMessage, error) {
	var msg MessageStatusEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	if msInfo != nil {
		message.clickURL = makeClickURL(msg.Message.ID, msInfo.Host, msInfo.Organization, msInfo.Name, "")
	}

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

func (p *Plugin) handleMessageLoadedEvent(payload json.RawMessage, msInfo *PostalMailserverInfo) (*GotifyMessage, error) {
	var msg MessageLoadedEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	if msInfo != nil {
		message.clickURL = makeClickURL(msg.Message.ID, msInfo.Host, msInfo.Organization, msInfo.Name, "/activity")
	}

	// TODO implement later
	message.Title = "not implemented"
	message.Message = "unimplemented message load event message"

	return message, nil
}

func (p *Plugin) handleMessageBounceEvent(payload json.RawMessage, msInfo *PostalMailserverInfo) (*GotifyMessage, error) {
	var msg MessageBounceEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	if msInfo != nil {
		message.clickURL = makeClickURL(msg.OriginalMessage.ID, msInfo.Host, msInfo.Organization, msInfo.Name, "")
	}

	// TODO implement later
	message.Title = "not implemented"
	message.Message = "unimplemented bounce event message"

	return message, nil
}

func (p *Plugin) handleMessageClickEvent(payload json.RawMessage, msInfo *PostalMailserverInfo) (*GotifyMessage, error) {
	var msg MessageClickEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	if msInfo != nil {
		message.clickURL = makeClickURL(msg.Message.ID, msInfo.Host, msInfo.Organization, msInfo.Name, "/activity")
	}

	// TODO implement later
	message.Title = "not implemented"
	message.Message = "unimplemented click event message"

	return message, nil
}

func (p *Plugin) handleDNSErrorEvent(payload json.RawMessage) (*GotifyMessage, error) {
	var msg DNSErrorEvent
	if err := json.Unmarshal(payload, &msg); err != nil {
		return nil, err
	}

	message := &GotifyMessage{}
	// this message type cannot have a click url, since the payload does not contain a message ID

	// TODO implement later
	message.Title = "not implemented"

	return message, nil
}
