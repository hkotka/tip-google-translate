package main

import (
	translate "cloud.google.com/go/translate/apiv3"
	"context"
	"fmt"
	"google.golang.org/api/option"
	translatepb "google.golang.org/genproto/googleapis/cloud/translate/v3"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	gcCredsFile  = "credentials.json"
	gcTargetLang = "en-US"
	gcProject    = "my-gcloud-project"
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
	} else {
		c := gcCreateClient()
		translation := gcTranslateText(c, inputText)
		text = fmt.Sprintf("[{\"type\": \"text\", \"value\": \"Google Translate\"},{\"type\": \"text\", \"value\": \"%s\"}]", translation)
	}
	fmt.Printf(text)
}

func gcCreateClient() *translate.TranslationClient {
	d := contextTime()
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	c, err := translate.NewTranslationClient(ctx, option.WithCredentialsFile(gcCredsFile))
	if err != nil {
		log.Fatal(err)
	}

	return c
}

func gcTranslateText(c *translate.TranslationClient, text []string) string {
	d := contextTime()
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()
	req := &translatepb.TranslateTextRequest{
		Contents:           text,
		TargetLanguageCode: gcTargetLang,
		Parent:             gcProject,
	}

	resp, err := c.TranslateText(ctx, req)
	if err != nil {
		log.Fatal(err)
	}
	return resp.Translations[0].DetectedLanguageCode + " - " + resp.Translations[0].TranslatedText
}

func contextTime() time.Time {
	t := time.Now().Add(15 * time.Second)
	return t
}
