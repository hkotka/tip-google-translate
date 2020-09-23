package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"time"

	translate "cloud.google.com/go/translate/apiv3"
	"google.golang.org/api/option"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
)

const (
	gcCredentialsFile = "credentials.json"
	gcTargetLanguage  = "en-US"
	gcProject         = "projects/my-gcloud-project"
	ctxTimeout        = 15
)

func main() {
	var text string
	var inputText = make([]string, 0)
	inputText = append(inputText, os.Args[1])
	strInputText := strings.Trim(fmt.Sprintf("%s", inputText), "[]")

	// Don't try to translate links
	isHTTP := regexp.MustCompile(`^https?://`)
	if isHTTP.MatchString(strInputText) {
		return
	}
	ctx, c := gcCreateClient()
	translation := gcTranslateText(ctx, c, inputText)
	text = fmt.Sprintf("[{\"type\": \"text\", \"value\": \"Google Translate\"},{\"type\": \"text\", \"value\": \"%s\"}]", translation)
	fmt.Printf(text)
}

func gcCreateClient() (context.Context, *translate.TranslationClient) {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx, option.WithCredentialsFile(gcCredentialsFile))
	if err != nil {
		log.Fatal(err)
	}

	return ctx, c
}

func gcTranslateText(ctx context.Context, c *translate.TranslationClient, text []string) string {
	ctx, cancel := context.WithDeadline(ctx, contextTimeout())
	defer cancel()
	req := &translatepb.TranslateTextRequest{
		Contents:           text,
		TargetLanguageCode: gcTargetLanguage,
		Parent:             gcProject,
	}

	resp, err := c.TranslateText(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Translations[0].DetectedLanguageCode + " - " + resp.Translations[0].TranslatedText
}

func contextTimeout() time.Time {
	t := time.Now().Add(ctxTimeout * time.Second)
	return t
}
