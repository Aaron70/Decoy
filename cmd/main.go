package main

import (
	"fmt"

	"github.com/aaron70/decoy"
	"github.com/aaron70/decoy/internal/utils"
)

func main() {
	d, err := decoy.NewDecoyWithSeed(0)
	utils.PanicErr(err)

	// var res map[string]any
	// err = json.Unmarshal([]byte(tmpl), &res)
	// utils.PanicErr(err)

	n := 10000
	t := 0
	f := 0

	for range n {
		tmpl, err := d.ParseTemplateString(`{{ Probability 0.9999 }}`)
		utils.PanicErr(err)

		if tmpl == "true" {
			t += 1
		} else {
			f += 1
		}
	}

	fmt.Printf("T: %f%%\n", (float64(t) / float64(n)) * 100.0)
	fmt.Printf("F: %f%%\n", (float64(f) / float64(n)) * 100.0)

}
