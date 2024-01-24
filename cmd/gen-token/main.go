package main

import (
	"fmt"
	"github.com/kazhuravlev/example-pow-guard/pkg/hashcash"
)

func main() {
	h := hashcash.NewStd()
	hash, err := h.Mint("xxx")
	if err != nil {
		panic("the sky is falling: " + err.Error())
	}

	fmt.Println(hash)
}
