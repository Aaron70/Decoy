package decoy

import "github.com/aaron70/decoy/internal/random"

var Default = must(NewDecoyWithSeed(0))

func must[T any](v T, err error) T {
	if err != nil {
		panic(err)
	}
	return v
}

func RandomChoice[T any](decoy *Decoy, choices ...T) T {
	return random.RandomChoice(decoy.Rand, choices...)
}

