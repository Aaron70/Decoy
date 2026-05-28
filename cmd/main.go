package main

import (
	"fmt"

	"github.com/aaron70/decoy"
	"github.com/aaron70/decoy/internal/utils"
)


func main() {
	d, err := decoy.NewDecoyWithSeed(10)
	utils.Panic(err)

	fmt.Printf("Random text: %q\n", d.RandomText(10))
	fmt.Printf("Incremental 1: %d\n", d.NextIncrementalInt("id", 1, 1))
	fmt.Printf("Incremental 2: %d\n", d.NextIncrementalInt("id", 1, 1))
	fmt.Printf("Incremental 3: %d\n", d.NextIncrementalInt("id", 1, 1))
	fmt.Printf("Incremental 4: %d\n", d.NextIncrementalInt("id", 1, 1))
	fmt.Printf("Template: %q\n", utils.Must(d.ParseTemplateString("Random Text: {{ RandomText 12 }}, Incremental: {{ NextIncrementalInt \"id\" 0 1 }}")))
}
