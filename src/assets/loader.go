package assets

import (
	"io/ioutil"
)

const basePath = "assets"

const (
    IconUnread = "gmail-red.ico"
    IconAllRead = "gmail-blue.ico"
)

func Get(s string) []byte {
    b, err := ioutil.ReadFile(basePath + "/" + s)
    if err != nil {
        panic(err)
    }
    return b
}