package decoy

import (
	"math/rand/v2"
	"sync/atomic"
	"time"

	"github.com/aaron70/decoy/internal/random"
)


type Decoy struct {
	Rand *rand.Rand
	randomText *random.RandomText
	intIncrementals map[string]*atomic.Int64
}

func NewDecoy(source rand.Source) (*Decoy, error) {
	rand := rand.New(source)
	d := &Decoy{
		Rand: rand,
		randomText: random.NewRandomText(rand, 6),
		intIncrementals: make(map[string]*atomic.Int64),
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

func (d *Decoy) NextIncrementalInt(id string, start, step int) int64 {
	incremental, exists := d.intIncrementals[id]
	if !exists {
		incremental = &atomic.Int64{}
		incremental.Store(int64(start))
		d.intIncrementals[id] = incremental
	} else {
		incremental.Add(int64(step))
	}
	return incremental.Load()
}

func (d *Decoy) RandomChoiceAny(choices ...any) any {
	return RandomChoice(d, choices...)
}
