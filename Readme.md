# Convertee

Converty twitch bot to listen in a channel and do forex conversions.

Listens for `convert [number] [currency_code] to [currency_code]`
e.g. `convert 60 usd to zar`

Also listens for `R[value]` and `$[value]` to show conversions between these two currencies.

## Building and Running

To start the bot and have it litsen in a twitch channel just pass the twitch username as the first arg.

```
go build
./convertee-twitch-bot twitchusername
```

**Expects a secrets.yaml file to be in the same directory as the executable.**

`./secrets.yaml`
```yaml
twitch_username: username
twitch_oauth: oauth:tokeotokentoenontek
fixer_api_key: asdasdasdasdadasdasdasdsd
```

### Twitch

Sign up for a normal twitch account.

You can get your twitch oauth token by visiting: https://twitchapps.com/tmi/
More info at: https://dev.twitch.tv/docs/irc

### Fixer

Sign up for a free fixer.io account
This package just uses the free API endpoint and that does up to 1000 requests per month.

The API key is on the dashboard.

### Deploying

To build for linux and upload to server
```
env GOOS=linux GOARCH=amd64 go build
scp ./convertee-twitch-bot deploy@serverip:/home/deploy
ssh deloy@serverip
vi secrets.yaml
./convertee-twitch-bot channel
```

## TODO
Tests