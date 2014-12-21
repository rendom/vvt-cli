package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

type Paste struct {
	Code      string `json:"code"`
	Encrypted string `json:"encrypted"`
	Language  string `json:"language"`
	DeleteAt  string `json:"delete_at"`
	CreatedAt string `json:"created_at"`
	Slug      string `json:"slug`
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

func postPaste(content string) Paste {
	data := Paste{content, "0", "", "", "", ""}
	body, err := json.Marshal(data)

	if err != nil {
		panic(err)
	}

	resp, err := http.Post("https://vvt.nu/save.json", "application/json", bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	body, err = ioutil.ReadAll(resp.Body)

	if err != nil {
		panic(err)
	}

	paste := decodeJson(body)
	return paste
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
	args := flag.Args()
	flag.Parse()

	if len(args) > 0 {
		if getSlug != nil {
			fmt.Printf(getPaste(*getSlug))
		}
	} else {
		bytes, err := ioutil.ReadAll(os.Stdin)
		if err != nil {
			panic(err)
		}
		content := string(bytes)
		paste := postPaste(content)
		fmt.Printf("https://vvt.nu/" + paste.Slug + "\n")
	}

}
