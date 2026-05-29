package decoy

import (
	"bytes"
	"text/template"
	"io"
	"math/rand/v2"
	"os"
	"sync/atomic"
	"time"

	"github.com/aaron70/decoy/internal/random"
)

type Decoy struct {
	Rand            *rand.Rand
	randomText      *random.RandomText
	intIncrementals map[string]*atomic.Int64
}

func NewDecoy(source rand.Source) (*Decoy, error) {
	rand := rand.New(source)
	d := &Decoy{
		Rand:            rand,
		randomText:      random.NewRandomText(rand, 6),
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

type parseTemplateOption func(*parseTemplateConfig) error

type parseTemplateConfig struct {
	Name  string
	Funcs template.FuncMap
	Data  any
}

func WithTemplateNamed(name string) parseTemplateOption {
	return func(ptc *parseTemplateConfig) error {
		ptc.Name = name
		return nil
	}
}

func WithData(data any) parseTemplateOption {
	return func(ptc *parseTemplateConfig) error {
		ptc.Data = data
		return nil
	}
}

func WithFuncMap(funcs template.FuncMap) parseTemplateOption {
	return func(ptc *parseTemplateConfig) error {
		ptc.Funcs = funcs
		return nil
	}
}

func (d *Decoy) DefaultTemplateFuncMaps() template.FuncMap {
	return template.FuncMap{
		"RandomInt":      d.RandomInt,
		"RandomChoice":   d.RandomChoiceAny,
		"RandomText":     d.RandomText,
		"NextIncrementalInt": d.NextIncrementalInt,
		"EnvVariable":    os.Getenv,
	}
}

func (d *Decoy) ParseTemplate(w io.Writer, tmpl string, options ...parseTemplateOption) error {
	config := &parseTemplateConfig{
		Funcs: d.DefaultTemplateFuncMaps(),
	}

	for _, option := range options {
		if err := option(config); err != nil {
			return err
		}
	}

	t, err := template.New(config.Name).
		Funcs(config.Funcs).
		Parse(tmpl)
	if err != nil {
		return err
	}
	return t.Execute(w, config.Data)
}

func (d *Decoy) ParseTemplateString(tmpl string, options ...parseTemplateOption) (string, error) {
	str := bytes.NewBufferString("")
	err := d.ParseTemplate(str, tmpl, options...)
	if err != nil {
		return "", err
	}
	return str.String(), nil
}
