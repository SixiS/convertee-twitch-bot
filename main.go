package main

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"convertee-twitch-bot/fixer"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/patrickmn/go-cache"
)

var mainCache *cache.Cache
var client *twitch.Client

func main() {
	client = twitch.NewClient("Convertee", "")
	mainCache = cache.New(30*time.Minute, 10*time.Minute)

	convertRegex := regexp.MustCompile(`(?i)convert ([\d\.]+) (\w+) to (\w+)`)
	randRegex := regexp.MustCompile(`R([\d\s\.]+)`)
	dollarRegex := regexp.MustCompile(`\$([\d\s\.]+)`)

	var amount float64
	var err error
	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if matched := convertRegex.FindAllStringSubmatch(message.Message, -1); len(matched) > 0 {
			amount, err = convertResponseHandler(matched[0][2], matched[0][3], matched[0][1])
			errorOrMessage(message.Channel, err, fmt.Sprintf("%s %s IS %.2f %s", matched[0][1], matched[0][2], amount, matched[0][3]))
		} else if matched := randRegex.FindAllStringSubmatch(message.Message, -1); len(matched) > 0 {
			amount, err = convertResponseHandler("ZAR", "USD", matched[0][1])
			errorOrMessage(message.Channel, err, fmt.Sprintf("R%s = $%.2f", matched[0][1], amount))
		} else if matched := dollarRegex.FindAllStringSubmatch(message.Message, -1); len(matched) > 0 {
			amount, err = convertResponseHandler("USD", "ZAR", matched[0][1])
			errorOrMessage(message.Channel, err, fmt.Sprintf("$%s = R%.2f", matched[0][1], amount))
		}
	})

	client.Join("Convertee")
	client.Say("Convertee", "/color SeaGreen")

	clientErr := client.Connect()
	if clientErr != nil {
		panic(clientErr)
	}
}

func errorOrMessage(channel string, err error, responseMessage string) {
	if err != nil {
		client.Say(channel, "/me || "+err.Error())
	} else {
		client.Say(channel, "/me || "+responseMessage)
	}
}

func convertResponseHandler(fromCurrency string, toCurrency string, amountAsString string) (float64, error) {
	var amount float64
	var err error

	if amount, err = strconv.ParseFloat(amountAsString, 64); err != nil {
		return 0, fmt.Errorf("Unable to convert %s %s to %s - %s is not a number", amountAsString, fromCurrency, toCurrency, amountAsString)
	}

	convertedAmount, err := fixer.Convert(mainCache, fromCurrency, toCurrency, amount)
	if err != nil {
		return 0, fmt.Errorf("Unable to convert %f %s to %s - conversion broke", amount, fromCurrency, toCurrency)
	}

	return convertedAmount, nil
}
