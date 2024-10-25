# Gotify Webhooks Plugin

### Installation

Just download the latest `.so` file for your architecture from the [package registry](https://git.leon.wtf/leon/gotify-postal-webhooks-plugin/-/packages) or build it yourself with `make build` (required Go and Docker). This uses Gotify's build tools to build against the latest version. The `.so` files are compiled to `build/gotify-postal-webhooks*.so`.

Then simply move the `.so` file to the Gotify plugin directory and restart Gotify.

### Usage

Activate the Plugin, then go to the plugin's details panel to retrieve the **Webhook URL**. You can also set you Postal instance details here. If configured, clicking messages redirects you to the Postal message dashboard.

The parsed payload is sent to the automatically created "Postal Webhooks" application channel along with the senders IP address. This channel can be renamed
