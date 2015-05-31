package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/rendom/vvt-cli/Godeps/_workspace/src/code.google.com/p/gopass"
	"github.com/rendom/vvt-cli/Godeps/_workspace/src/golang.org/x/crypto/ssh/terminal"
)

type Paste struct {
	Code      string `json:"code"`
	Encrypted bool   `json:"encrypted"`
	Language  string `json:"language"`
	Slug      string `json:"slug"`
}

// Get paste from api by passing slug
// Will prompt password if paste is encrypted
func GetPaste(slug string) string {
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
	paste := decodeJSON(body)

	// Move this logix out of GetPaste..
	if paste.Encrypted == true {
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

// Post paste trought api requires string and retunrs Paste object
func PostPaste(content string) Paste {
	data := Paste{Code: content, Encrypted: false}
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

	paste := decodeJSON(body)
	return paste
}

// decode vvt's json response to struct
func decodeJSON(s []byte) Paste {
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
			paste := PostPaste(string(bytes))
			fmt.Printf("https://vvt.nu/" + paste.Slug + "\n")
		} else {
			fmt.Println("No piped data")
			flag.Usage()
			return
		}
	case 1:
		content := GetPaste(args[0])
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
