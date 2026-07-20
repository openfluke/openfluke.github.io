package main

import (
	"fmt"

	"github.com/openfluke/welvet/stub/templates"
)

func main() {
	p := templates.ChatML.BuildPrompt(nil, "You are helpful.", "Say hi")
	fmt.Println(p)
}
