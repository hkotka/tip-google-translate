#!/bin/zsh
LANG="${TIP_LANG:-en}"
go build -ldflags="-s -w -X main.targetLanguage=$LANG"
cp tip-google-translate ~/Library/Application\ Scripts/tanin.tip/provider.script