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
	gcCredentialsFile = "credentials.json"
	gcTargetLanguage  = "en-US"
	gcProject         = "projects/my-gcloud-project"
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
		c, ctx := gcCreateClient()
		translation := gcTranslateText(c, ctx, inputText)
		text = fmt.Sprintf("[{\"type\": \"text\", \"value\": \"Google Translate\"},{\"type\": \"text\", \"value\": \"%s\"}]", translation)
	}
	fmt.Printf(text)
}

func gcCreateClient() (*translate.TranslationClient, context.Context) {
	ctx := context.Background()
	c, err := translate.NewTranslationClient(ctx, option.WithCredentialsFile(gcCredentialsFile))
	if err != nil {
		log.Fatal(err)
	}

	return c, ctx
}

func gcTranslateText(c *translate.TranslationClient, ctx context.Context, text []string) string {
	ctx, cancel := context.WithDeadline(ctx, contextTime())
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

func contextTime() time.Time {
	t := time.Now().Add(15 * time.Second)
	return t
}
