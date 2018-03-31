#!/bin/bash

rm -rf dist/*
mkdir -p dist/notifier-win/config
cp -r assets dist/notifier-win/assets
cp -r config/config.json.dist dist/notifier-win/config/config.json

GOOS=windows GOARCH=386 go build -ldflags -H=windowsgui -o dist/notifier-win/email-notifier.exe