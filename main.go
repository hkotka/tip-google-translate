package main

import (
	"cloud.google.com/go/translate"
	"context"
	_ "embed"
	"fmt"
	"golang.org/x/text/language"
	"google.golang.org/api/option"
	"os"
	"regexp"
	"strings"
	"time"
)

//go:embed credentials.json
var gcCredentialsJson []byte

const (
	gcTargetLanguage = "en"
	ctxTimeout       = 15
)

func main() {

	// Don't try to translate links
	isHTTP := regexp.MustCompile(`^https?://`)
	if isHTTP.MatchString(strings.Trim(fmt.Sprintf("%s", os.Args[1]), "[]")) {
		return
	}

	translation, err := translateText(gcTargetLanguage, os.Args[1])
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf(fmt.Sprintf("[{\"type\": \"text\", \"value\": \"Google Translate\"},{\"type\": \"text\", \"value\": \"%s\"}]", translation))
}

func translateText(targetLanguage, text string) (string, error) {
	ctx := context.Background()
	ctx, cancel := context.WithDeadline(ctx, contextTimeout())
	defer cancel()

	lang, err := language.Parse(targetLanguage)
	if err != nil {
		return "", fmt.Errorf("language.Parse: %v", err)
	}

	client, err := translate.NewClient(ctx, option.WithCredentialsJSON(gcCredentialsJson))
	if err != nil {
		return "", err
	}
	defer client.Close()

	resp, err := client.Translate(ctx, []string{text}, lang, nil)
	if err != nil {
		return "", fmt.Errorf("translate: %v", err)
	}
	if len(resp) == 0 {
		return "", fmt.Errorf("translate returned empty response to text: %s", text)
	}
	return resp[0].Text, nil
}

func contextTimeout() time.Time {
	t := time.Now().Add(ctxTimeout * time.Second)
	return t
}
