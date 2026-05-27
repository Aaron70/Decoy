package decoy

import (
	"math/rand/v2"
	"time"

	"github.com/aaron70/decoy/internal/random"
)


type Decoy struct {
	Rand *rand.Rand
	randomText *random.RandomText
}

func NewDecoy(source rand.Source) (*Decoy, error) {
	rand := rand.New(source)
	d := &Decoy{
		Rand: rand,
		randomText: random.NewRandomText(rand, 6),
	}
	return d, nil
}

func NewDecoyWithSeed(seed uint64) (*Decoy, error) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	return NewDecoy(rand.NewPCG(seed, 1))
}

func (d *Decoy) RandomText(maxWords int) string {
	return d.randomText.RandomText(maxWords)
}

func (d *Decoy) RandomInt(min, max int) int {
	return random.RandomInt(d.Rand, min, max)
}

func (d *Decoy) RandomChoiceAny(choices ...any) any {
	return RandomChoice(d, choices...)
}
