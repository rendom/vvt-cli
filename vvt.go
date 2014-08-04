package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Paste struct {
	Code      string `json:"code"`
	Encrypted string `json:"encrypted"` // Should be boolean, API gives wrong output atm.
	Language  string `json:"language"`
	Slug      string `json:"slug"`
}

func getPaste(slug string) string {
	url := "https://vvt.nu/" + slug + ".json"
	result, err := http.Get(url)

	if err != nil {
		panic(err)
	}

	defer result.Body.Close()

	body, err := ioutil.ReadAll(result.Body)

	if err != nil {
		panic(err)
	}
	paste := decodeJson(body)
	return paste.Code
}

func decodeJson(s []byte) Paste {
	var paste Paste

	if err := json.Unmarshal(s, &paste); err != nil {
		fmt.Printf("Hmm.. problem")
		os.Exit(1)
	}
	return paste
}

func main() {
	getSlug := flag.String("get", "get", "get")
	flag.Parse()

	if getSlug != nil {
		fmt.Printf(getPaste(*getSlug))
	}
}
