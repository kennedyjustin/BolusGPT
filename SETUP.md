# Setup Guide

This guide assumes a moderate degree of technical knowledge (how to use a command line, spin up EC2 instances, etc.). I'll walk through how I decided to setup my own instance of BolusGPT, and it is upon the user to deviate from my choices when desired.

## Caution

1. The Dexcom API we use is not official, and was reversed engineered by the [pydexcom project](https://github.com/gagebenne/pydexcom). Assume this API could change at any time, or that Dexcom could take action to prevent usage of their API in this manner.
1. By self-hosting you are making yourself responsible for the availability of the server. There is no one to call if it goes down for some reason.
1. Due to the above, this utility should be treated as a quality of life improvement, and not depended upon. Always have a backup way to calculate bolus doses.

## Decide where to self-host

Requirements:
1. Local (and ideally non-ephemeral) file storage available.
1. HTTPS->HTTP Reverse-Proxy available - OpenAI requires the OpenAPI spec to be hosted on a domain over HTTPS. Our server exposes an HTTP API.
1. A domain name for (2).

I personally chose to host on a [free-tier EC2 instance](https://aws.amazon.com/free/) behind an [nginx reverse proxy](https://docs.nginx.com/nginx/admin-guide/web-server/reverse-proxy/). I purchased my domain through [Namecheap](https://namecheap.com/).

## Dependencies

Get these installed on your EC2 Instance, and open it up to TCP ports 22 (SSH) and 443 (HTTPS) via Security Groups.

1. The [Go](https://go.dev/) language
1. [Nginx](https://docs.nginx.com)
1. Some way to keep the server running (I use [tmux](https://github.com/tmux/tmux/wiki))

## DNS

In your domain provider's DNS settings, create an `A` record pointing to the EC2 Instance's public IP address. Update [`openapi.yaml`](./openapi.yaml)'s `servers` field with your domain name (include `https://`.)

## TLS

Use [certbot](https://certbot.eff.org/instructions?ws=nginx&os=ubuntufocal) to get Let's Encrypt TLS certificates set up. In `nginx` create a reverse proxy from your domain to `localhost:8080`. It should look something like this

```
server {
    listen       443 ssl;
    listen       [::]:443 ssl;
    http2        on;
    server_name  api.bolusgpt.com;
    root         /usr/share/nginx/html;

    ssl_certificate "/etc/letsencrypt/live/bolusgpt.com/fullchain.pem";
    ssl_certificate_key "/etc/letsencrypt/live/bolusgpt.com/privkey.pem";
    ssl_session_cache shared:SSL:1m;
    ssl_session_timeout  10m;
    ssl_ciphers PROFILE=SYSTEM;
    ssl_prefer_server_ciphers on;

    location / {
        proxy_pass http://localhost:8080;
        proxy_pass_request_headers on;
    }
```

## Start the Server

Create a UUID that will be used as the API's Bearer token:

```
uuidgen
```

Run the API server (I do this in a `tmux` session):
```
DEXCOM_USERNAME="<username>" DEXCOM_PASSWORD="<password>" BEARER_TOKEN="<token>" TZ="America/New_York" sudo -E go run .
```

Try using the API. Here are a few examples:

```
# Get user settings (should all be default values to start with):
curl -X GET -H "Authorization: Bearer <token>" https://<domain>/me

# Update a setting:
curl -X PATCH -H "Authorization: Bearer <token>" https://<domain>/me -d '{"target_blood_glucose_level_in_mg_dl": 100, "insulin_to_carb_ratio": 6, "insulin_sensitivity_factor": 20}'

# Try to calculate a dose:
curl -X POST -H "Authorization: Bearer <token>" https://<domain>/dose -d '{"total_grams_of_carbs": 20}'
```

## Create the Custom GPT

In your own OpenAI account, create a new CustomGPT. Provide the following configuration:
- **Name**: BolusGPT
- **Description**: Whatever you want
- **Picture**: Whatever you want.
- **Instructions**: Copy and paste from [`SPEC.md`](./SPEC.md).
- **Conversation starters**:
   - "Onboard"
   - "Get Dose"
   - "Get Settings"
   - "Confirm Dose"
- **Actions**: Copy [`openapi.yaml`](./openapi.yaml). Make sure to update with your own server domain.

Finally, create the Custom GPT. I keep the Share settings to "Only me"

## Test it out

Finally, try onboarding and testing out the dose algorithm via text or voice.