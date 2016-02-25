package main

import (
	"fmt"
	"github.com/mmontagnino/gbparser"
)

func main() {
	fmt.Printf(gbparser.Parse("sequence.gb"))
}
