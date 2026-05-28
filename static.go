package decoy

import (
	"github.com/aaron70/decoy/internal/random"
	"github.com/aaron70/decoy/internal/utils"
)

var Default = utils.Must(NewDecoyWithSeed(0))


func RandomChoice[T any](decoy *Decoy, choices ...T) T {
	return random.RandomChoice(decoy.Rand, choices...)
}

