# Convertee

Converty twitch bot to listen in channels and do forex conversions.

```
go build
./convertee-twitch-bot channel
```

To build for linux and upload to server
```
env GOOS=linux GOARCH=amd64 go build
scp ./convertee-twitch-bot deploy@serverip:/home/deploy
ssh deloy@serverip
./convertee-twitch-bot channel
```