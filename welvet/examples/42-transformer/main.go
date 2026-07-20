package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/tokenizer"
	"github.com/openfluke/welvet/model/transformer"
)

func main() {
	m, err := transformer.LoadEntity("model.entity")
	if err != nil {
		fmt.Println(err)
		return
	}
	tok, err := tokenizer.LoadTokenizer("tokenizer.json")
	if err != nil {
		panic(err)
	}
	if err := m.ApplyExec(transformer.ProfileSIMDMultiCore()); err != nil {
		panic(err)
	}
	text, _, err := m.Generate(
		tok.Encode, tok.Decode, nil,
		"You are helpful.", "Say hi.",
		transformer.GenOptions{MaxTokens: 32, Silent: true},
	)
	fmt.Println(text, err)
}
