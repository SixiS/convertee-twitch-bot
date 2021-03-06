package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"convertee-twitch-bot/fixer"
	"convertee-twitch-bot/googletranslate"

	"github.com/gempir/go-twitch-irc/v2"
	"github.com/patrickmn/go-cache"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

var mainCache *cache.Cache
var client *twitch.Client
var googleTranslateClient *googletranslate.TranslateClient

func main() {
	argsWithoutProg := os.Args[1:]
	if len(argsWithoutProg) == 0 {
		panic("Channel must be provided as command line arg.")
	}
	channel := argsWithoutProg[0]

	viper.SetConfigName("secrets")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error loading config file: %s", err))
	}

	googleTranslateClient = googletranslate.NewTranslateClient()

	client = twitch.NewClient(viper.GetString("twitch_username"), viper.GetString("twitch_oauth"))
	mainCache = cache.New(30*time.Minute, 10*time.Minute)

	convertRegex := regexp.MustCompile(`(?i)convert ([\d\.]+) (\w+) to (\w+)`)
	randRegex := regexp.MustCompile(`\bR([\d\.]+)\b`)
	dollarRegex := regexp.MustCompile(`\$([\d\.]+)\b`)
	translateRegex := regexp.MustCompile(`(?i)\!translate (.+)`)

	client.OnConnect(func() {
		fmt.Println("Connected and listening!")
		twitchColour := viper.GetString("twitch_colour")
		if len(twitchColour) > 0 {
			fmt.Println("Setting colour to", twitchColour)
			client.Say(channel, fmt.Sprintf("/color %s", twitchColour))
		}
	})

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		if matched := convertRegex.FindAllStringSubmatch(message.Message, -1); len(matched) > 0 {
			amount, err := convertAmount(matched[0][2], matched[0][3], matched[0][1])
			errorOrMessage(message.Channel, err, fmt.Sprintf("%s %s IS %.2f %s", matched[0][1], matched[0][2], amount, matched[0][3]))
		} else if matched := randRegex.FindAllStringSubmatch(message.Message, -1); len(matched) > 0 {
			amount, err := convertAmount("ZAR", "USD", matched[0][1])
			errorOrMessage(message.Channel, err, fmt.Sprintf("R%s = $%.2f", matched[0][1], amount))
		} else if matched := dollarRegex.FindAllStringSubmatch(message.Message, -1); len(matched) > 0 {
			amount, err := convertAmount("USD", "ZAR", matched[0][1])
			errorOrMessage(message.Channel, err, fmt.Sprintf("$%s = R%.2f", matched[0][1], amount))
		} else if matched := translateRegex.FindAllStringSubmatch(message.Message, -1); len(matched) > 0 {
			translatedText, err := googleTranslateClient.Translate(matched[0][1], language.English)
			errorOrMessage(message.Channel, err, fmt.Sprintf("%q in english is %q", matched[0][1], translatedText))
		}
	})

	client.Join(channel)

	for true {
		fmt.Println("Connecting to #" + channel)
		clientErr := client.Connect()
		if clientErr != nil {
			fmt.Println("Error with connection to twitch... sleeping 5 seconds and retrying.")
			time.Sleep(2 * time.Second)
		}
	}
}

func errorOrMessage(channel string, err error, responseMessage string) {
	if err != nil {
		client.Say(channel, "/me : "+err.Error())
	} else {
		client.Say(channel, "/me : "+responseMessage)
	}
}

func convertAmount(fromCurrency string, toCurrency string, amountAsString string) (float64, error) {
	var amount float64
	var err error

	if amount, err = strconv.ParseFloat(amountAsString, 64); err != nil {
		return 0, fmt.Errorf("Unable to convert %q %s to %s - %q is not a number", amountAsString, fromCurrency, toCurrency, amountAsString)
	}

	convertedAmount, err := fixer.Convert(mainCache, fromCurrency, toCurrency, amount)
	if err != nil {
		return 0, fmt.Errorf("Unable to convert %q %s to %s - %s", amountAsString, fromCurrency, toCurrency, err.Error())
	}

	return convertedAmount, nil
}
