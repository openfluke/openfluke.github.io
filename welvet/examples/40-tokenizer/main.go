package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/tokenizer"
)

func main() {
	tok, err := tokenizer.LoadTokenizer("tokenizer.json")
	if err != nil {
		fmt.Println(err)
		return
	}
	ids := tok.Encode("hello welvet", true)
	fmt.Println(ids, tok.Decode(ids, true))
}
