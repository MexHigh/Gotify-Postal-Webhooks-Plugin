package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

const (
	EmojiCheckMark   = "\xE2\x9C\x85"
	EmojiExclamMark  = "\xE2\x9D\x97"
	EmojiWarningSign = "\xE2\x9A\xA0"
)

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
		message.Title = EmojiCheckMark + " Message delivered successfully"
	case WebhookMessageEventMessageDelayed:
		message.Title = EmojiWarningSign + " Message delivery delayed"
	case WebhookMessageEventMessageDeliveryFailed:
		message.Title = EmojiExclamMark + " Message delivery failed"
	case WebhookMessageEventMessageHeld:
		message.Title = EmojiWarningSign + " Message delivery was held by Postal"
	default:
		return nil, errors.New("unknown event name '" + string(eventType) + "' occured in message status event handler")
	}

	message.Message += fmt.Sprintf("_From %s to %s: \"%s\"_\n\n", msg.Message.From, msg.Message.To, msg.Message.Subject)
	message.Message += msg.Details + "\n\n"
	message.Message += "---\n\n"
	if msg.Time != 0.0 {
		message.Message += fmt.Sprintf("**Delivery time:** %f seconds\n\n", msg.Time)
	} else {
		message.Message += "**Delivery time:** instant\n\n"
	}
	message.Message += fmt.Sprintf("**Sent with SSL/TLS:** %t\n\n", msg.SentWithSSL)
	if msg.Output != "" {
		message.Message += fmt.Sprintf("**Output:**\n\n```\n%s\n```", msg.Output)
	} else {
		message.Message += "**Output:** none"
	}

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
