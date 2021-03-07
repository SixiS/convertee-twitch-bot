package googletranslate

import (
	"context"
	"fmt"

	"cloud.google.com/go/translate"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
)

// NewTranslateClient returns an initialised google translate service api client
func NewTranslateClient() *TranslateClient {
	ctx := context.Background()
	translateService, err := translate.NewClient(ctx, option.WithAPIKey(viper.GetString("google_translate_api_key")))
	if err != nil {
		panic(err)
	}
	return &TranslateClient{
		context: ctx,
		client:  translateService,
	}
}

// TranslateClient is a struct holding an initialised google translation api service
type TranslateClient struct {
	context context.Context
	client  *translate.Client
}

// Translate uses the google translate API to translate the source text to the target language
func (t *TranslateClient) Translate(text string, targetLanguage language.Tag) (string, error) {
	resp, err := t.client.Translate(t.context, []string{text}, targetLanguage, nil)
	if err != nil {
		return "", fmt.Errorf("Error translating: %s", err)
	}
	return resp[0].Text, nil
}
