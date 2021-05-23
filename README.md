# DDNS

![Test](https://github.com/skibish/ddns/workflows/run%20tests/badge.svg)
![Release](https://github.com/skibish/ddns/workflows/release/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/skibish/ddns)](https://goreportcard.com/report/github.com/skibish/ddns)

Personal DDNS client with [Digital Ocean Networking](https://www.digitalocean.com/products/networking/) DNS as backend.

*[Read about it in the Blog](https://sergeykibish.com/blog/ddns-v4)*

## Motivation

There are services like [DynDNS](http://dyn.com/dns/), [No-IP](http://www.noip.com/) to access PCs remotely.
But do we need them?
This is your own DDNS solution which works for free (thanks to [Digital Ocean Networking](https://www.digitalocean.com/products/networking/) DNS).

## What is DDNS

*From [Wikipedia](https://en.wikipedia.org/wiki/Dynamic_DNS)*
> Dynamic DNS (DDNS) is a method of automatically updating a name server in the Domain Name System (DNS), often in real time, with the active DDNS configuration of its configured hostnames, addresses or other information.

## Installation

Download binary from [releases](https://github.com/skibish/ddns/releases).

And start it as:

```bash
ddns
```

Or you can download [Docker image](https://hub.docker.com/r/skibish/ddns) and use it:

```bash
docker run \
  -v /path/to/config.yml:/config/ddns.yml \
  skibish/ddns -conf-file /config/ddns.yml
```

## Documentation

Run `ddns -h`, to see help.
It will output:

```text
Usage of ./ddns:
  -conf-file string
        Location of the configuration file. If not provided, searches current directory, then $HOME for ddns.yml file
  -ver
        Show version
```

You need to setup your domain in Digital Ocean Networks panel.

In your domain name provider configuration point domain to Digital Ocean NS records.

*Refer to: [How To Point to DigitalOcean Nameservers From Common Domain Registrars](https://www.digitalocean.com/community/tutorials/how-to-point-to-digitalocean-nameservers-from-common-domain-registrars)*

Configuration file `ddns.yml`:

```yaml
# DDNS configuration file.

# Mandatory, DigitalOcean API token.
# It can be also set using environment variable DDNS_TOKEN.
token: ""

# By default, IP check occurs every 5 minutes.
# It can be also set using environment variable DDNS_CHECKPERIOD.
checkPeriod: "5m"

# By default, timeout to external resources is set to 10 seconds.
# It can be also set using environment variable DDNS_REQUESTTIMEOUT.
requestTimeout: "10s"

# By default, IPv6 address is not requested.
# IPv6 address can be forced by setting it to `true`.
# It can be also set using environment variable DDNS_IPV6.
ipv6: false

# List of domains and their records to update.
domains:
  example.com:
  # More details about the fields can be found here:
  # https://developers.digitalocean.com/documentation/v2/#create-a-new-domain-record
  - type: "A"
    name: "www"
  - type: "TXT"
    name: "demo"

    # By default, is set to "{{.IP}}" (key .IP is reserved).
    # Supports Go template engine.
    # Additional keys can be set in "params" block below.
    data: "My IP is {{.IP}} and I am {{.mood}}"

    # By default, 1800 seconds (5 minutes).
    ttl: 1800

# By default, params is empty.
params:
  mood: "cool"

# By default, notifications is empty.
notifications:

  # Gotify (https://gotify.net)
- type: "gotify"
  app_url: "https://gotify.example.com"
  app_token: ""
  title: "DDNS" 

  # SMTP
- type: "smtp"
  user: "foo@bar.com"
  password: "1234"
  host: "localhost"
  port: "468"
  from: "bar@foo.com"
  to: "foo@foo.com"
  subject: "My DDNS sending me a message"

  # Telegram (https://telegram.org)
- type: "telegram"
  token: "telegram bot token"
  chat_id: "1234"
```
