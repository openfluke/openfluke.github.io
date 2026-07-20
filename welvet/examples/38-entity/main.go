package main

import (
	"fmt"

	"github.com/openfluke/welvet/model/entity"
)

func main() {
	info, err := entity.Inspect("model.entity")
	if err != nil {
		fmt.Println(err)
		return
	}
	ef, err := entity.Open("model.entity")
	if err != nil {
		panic(err)
	}
	defer ef.Close()
	fmt.Println(info, ef.HasTokenizerBlob())
}
