package main

import (
	"fmt"

	"github.com/aaron70/decoy"
)

func Panic(err error) {
	if err != nil {
		panic(err)
	}
}

func main() {
	d, err := decoy.NewDecoyWithSeed(10)
	Panic(err)

	fmt.Printf("Random text: %q\n", d.RandomText(10))
}
