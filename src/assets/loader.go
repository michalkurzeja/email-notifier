package assets

import (
	"path/filepath"
	"os"
	"io/ioutil"
)

const basePath = "assets"

const (
    IconUnread string = "gmail-red.ico"
    IconAllRead string = "gmail-blue.ico"
)

func Get(asset string) []byte {
    b, err := ioutil.ReadFile(GetPath(asset))
    if err != nil {
        panic(err)
    }
    return b
}

func GetPath(asset string) string {
    return basePath + "/" + asset
}

func GetAbsolutePath(asset string) string {
    return getCurrentPath() + "/" + GetPath(asset)
}

func getCurrentPath() string {
    ex, err := os.Executable()
    if err != nil {
        panic(err)
    }
    exPath := filepath.Dir(ex)

    return exPath
}