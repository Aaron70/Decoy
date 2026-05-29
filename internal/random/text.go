package random

import (
	"maps"
	"math/rand/v2"
	"slices"
	"strings"
_ "embed"
)

//go:embed corpus.txt
var DefaultRandomTextCorpus string

type RandomText struct {
	Rand             *rand.Rand
	n                int
	ngrams           map[string][]string
}

func NewRandomText(rand *rand.Rand, n int) *RandomText {
	return &RandomText{
		n:      n,
		Rand:   rand,
		ngrams: make(map[string][]string),
	}
}

func (r *RandomText) SetNgrams(n int, ngrams map[string][]string) {
	r.n = n
	r.ngrams = ngrams
}

func (r *RandomText) NgramsFromString(text string) {
	if len(r.ngrams) != 0 {
		r.ngrams = make(map[string][]string)
	}
	words := strings.Fields(text)
	r.NgramsFromWords(words)
}

func (r *RandomText) NgramsFromWords(words []string) {
	if len(words) <= r.n {
		return
	}
	for i := 0; i < len(words)-r.n; i++ {
		key := strings.Join(words[i:i+r.n], " ")
		r.ngrams[key] = append(r.ngrams[key], words[i+r.n])
	}
}

func (r *RandomText) RandomText(maxWords int) string {
	if len(r.ngrams) == 0 {
		r.NgramsFromString(DefaultRandomTextCorpus)
	}
	keys := slices.Collect(maps.Keys(r.ngrams))
	words := strings.Fields(RandomChoice(r.Rand, keys...))
	text := ""

	for range maxWords {
		key := strings.Join(words[len(words)-r.n:], " ")
		successors, ok := r.ngrams[key]
		if !ok {
			break
		}
		words = append(words, RandomChoice(r.Rand, successors...))
		if text != "" {
			text += " "
		}
		text += words[len(words)-1]
	}

	return text
}
