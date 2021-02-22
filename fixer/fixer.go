package fixer

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/patrickmn/go-cache"
)

type convertResponse struct {
	Success bool `json:"success"`
	Query   struct {
		From   string `json:"from"`
		To     string `json:"to"`
		Amount int    `json:"amount"`
	} `json:"query"`
	Info struct {
		Timestamp int     `json:"timestamp"`
		Rate      float64 `json:"rate"`
	} `json:"info"`
	Date   string  `json:"date"`
	Result float64 `json:"result"`
	Error  struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
}

type listResponse struct {
	Success   bool               `json:"success"`
	Timestamp int                `json:"timestamp"`
	Base      string             `json:"base"`
	Date      string             `json:"date"`
	Rates     map[string]float64 `json:"rates"`
	Error     struct {
		Code int    `json:"code"`
		Type string `json:"type"`
		Info string `json:"info"`
	} `json:"error"`
}

const fixerAccessKey string = ""

// Convert converts an amount using the Fixer List API
func Convert(cache *cache.Cache, fromCurrency string, toCurrency string, amount float64) (float64, error) {
	var listObject listResponse

	if x, found := cache.Get("fixerList"); found {
		listObject = *x.(*listResponse)
	} else {
		var err error
		listObject, err = fixerList()
		if err != nil {
			return 0, err
		}
		cache.SetDefault("fixerList", &listObject)
	}

	fromRate, ok := listObject.Rates[strings.ToUpper(fromCurrency)]
	if !listObject.Success || !ok {
		return 0, errors.New("From currency not supported")
	}

	toRate, ok := listObject.Rates[strings.ToUpper(toCurrency)]
	if !listObject.Success || !ok {
		return 0, errors.New("To currency not supported")
	}

	EURPrice := amount / fromRate

	return toRate * EURPrice, nil
}

func fixerList() (listResponse, error) {
	listRequestURL := fmt.Sprintf("http://data.fixer.io/api/latest?access_key=%s", fixerAccessKey)
	response, err := http.Get(listRequestURL)
	if err != nil {
		return listResponse{}, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	var listObject listResponse
	json.Unmarshal(responseData, &listObject)
	return listObject, nil
}

func fixerConvert(fromCurrency string, toCurrency string, amount float64) (float64, error) {
	conversionRequestURL := fmt.Sprintf("http://data.fixer.io/api/convert?access_key=%s&from=%s&to=%s&amount=%f", fixerAccessKey, fromCurrency, toCurrency, amount)
	response, err := http.Get(conversionRequestURL)
	if err != nil {
		return 0, err
	}
	responseData, err := ioutil.ReadAll(response.Body)
	var responseObject convertResponse
	json.Unmarshal(responseData, &responseObject)
	if responseObject.Success {
		return responseObject.Result, nil
	} else {
		return 0, errors.New("Failed")
	}
}
