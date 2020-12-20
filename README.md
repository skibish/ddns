# DDNS

![Test](https://github.com/skibish/ddns/workflows/test/badge.svg)
![Release](https://github.com/skibish/ddns/workflows/release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/skibish/ddns)](https://goreportcard.com/report/github.com/skibish/ddns)

Personal DDNS client with [Digital Ocean Networking](https://www.digitalocean.com/products/networking/) DNS as backend.

*[Read about it in the Blog](https://sergeykibish.com/blog/your-personal-ddns)*

## Motivation

We have services like [DynDNS](http://dyn.com/dns/), [No-IP](http://www.noip.com/) to access PCs remotely. But do we need them?
This project is your own DDNS solution and will work for free (thanks to [Digital Ocean Networking](https://www.digitalocean.com/products/networking/) DNS).

## What is DDNS

*From [Wikipedia](https://en.wikipedia.org/wiki/Dynamic_DNS)*
> Dynamic DNS (DDNS or DynDNS) is a method of automatically updating a name server in the Domain Name System (DNS), often in real time, with the active DDNS configuration of its configured hostnames, addresses or other information.

## Installation

Download binary from [releases](https://github.com/skibish/ddns/releases) to `/usr/local/bin/ddns`.

And start it as:

```bash
ddns
```

Or you can download [Docker image](https://hub.docker.com/r/skibish/ddns) and use it:

```bash
docker run \
  -v /path/to/config.yml:/config/.ddns.yml \
  skibish/ddns -conf-file /config/.ddns.yml
```

## Documentation

You can download binary for your OS from [releases page](https://github.com/skibish/ddns/releases).

> **ATTENTION!** Currently tested on Linux and macOS.

Run `ddns -h`, to see help. It will output:

```text
Usage of ./ddns:
  -check-period duration
      Check if IP has been changed period (default 5m0s)
  -conf-file string
      Location of the configuration file (default "$HOME/.ddns.yml")
  -req-timeout duration
      Request timeout to external resources (default 10s)
  -v  Show version and exit
```

**Configuration should be supplied.** By default it is read from `$HOME/.ddns.yml`.

You need to setup your domain in Digital Ocean Networks panel.

In your domain name provider configuration point domain to Digital Ocean NS records.

*Refer to: [How To Set Up a Host Name with DigitalOcean](https://www.digitalocean.com/community/tutorials/how-to-set-up-a-host-name-with-digitalocean)*

Configuration should be in the following format:

```yaml
token: "AMAZING TOKEN"                          # Digital Ocean token
domains:                                        # Domains to update
  - "example.com"
forceIPV6: true                                 # Use IPv6 address resolve (Default false and force IPv4)
records:                                        # Records of the domains to update
  - type: "A"                                   # Record type
    name: "www"                                 # Record name
  - type: "TXT"
    name: "demo"
    data: "My IP is {{.IP}} and I am {{.mood}}" # "data" key is optional. You can write here
                                                # what you want and reference values from "params".
                                                # Key "IP" is reserved.
params:
  mood: "cool"
notify:                                         # Optional notifiers
  smtp:
    read: below
  telegram:
    read: below
```

### Notifications

These notifications are based on [sirupsen/logrus hooks](https://github.com/sirupsen/logrus#hooks).
Add them to the configuration file as:

```yaml
# config part from the top
#...

notify:
  <name of notification>:
    # ...configuration
```

List of supported notifications:

#### SMTP

```yaml
smtp:
  user: "foo@bar.com"
  password: "1234"
  host: "localhost"
  port: "22"
  to: "foo@foo.com"
  subject: "My DDNS sending me a message"
  secure: true # Optional flag. Set it, if you will send emails with SSL
```

#### Telegram

```yaml
telegram:
  token: "telegram bot token"
  chat_id: "1234"
```
