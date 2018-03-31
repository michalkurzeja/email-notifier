#!/bin/bash

rm -rf dist/*
mkdir dist/notifier-win
cp -r assets dist/notifier-win/assets

GOOS=windows GOARCH=386 go build -ldflags -H=windowsgui -o dist/notifier-win/email-notifier.exe