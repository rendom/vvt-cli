package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"code.google.com/p/gopass"
	"golang.org/x/crypto/ssh/terminal"
)

type Paste struct {
	Id        int    `json:"id"`
	Code      string `json:"code"`
	Encrypted int    `json:"encrypted"`
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

	if paste.Encrypted == 1 {
		password, err := gopass.GetPass("This paste is encrypted, enter password: ")
		if err != nil {
			panic(err)
		}

		content, err := Decrypt(paste.Code, password)
		if err != nil {
			panic(err)
		}
		return content
	}

	return paste.Code
}

func postPaste(content string) Paste {
	data := Paste{0, content, 0, "", "", "", ""}
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
		panic(err)
		fmt.Printf("Hmm.. problem")
		os.Exit(1)
	}
	return paste
}

func main() {
	outputFile := flag.String("o", "", "Output file (defaul is stdout)")
	flag.Parse()

	args := flag.Args()
	switch len(args) {
	case 0:
		if !terminal.IsTerminal(0) {
			bytes, err := ioutil.ReadAll(os.Stdin)
			if err != nil {
				panic(err)
			}

			paste := postPaste(string(bytes))
			fmt.Printf("https://vvt.nu/" + paste.Slug + "\n")
		} else {
			fmt.Println("No piped data")
			flag.Usage()
			return
		}
	case 1:
		content := getPaste(args[0])
		if *outputFile != "" {
			file, err := os.Create(*outputFile)
			if err != nil {
				panic(err)
			}
			file.WriteString(content)
			fmt.Println(*outputFile + " created")
		} else {
			fmt.Print(content)
		}
	default:
		flag.Usage()
		os.Exit(-1)
	}
}
