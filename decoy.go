package decoy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/rand/v2"
	"os"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/aaron70/decoy/internal/random"
	"github.com/aaron70/goaty/validations"
)

type Decoy struct {
	Rand            *rand.Rand
	randomText      *random.RandomText
	intIncrementals map[string]*atomic.Int64
	templates       map[string]*template.Template
	m               sync.RWMutex
}

func NewDecoy(source rand.Source) (*Decoy, error) {
	rand := rand.New(source)
	d := &Decoy{
		Rand:            rand,
		randomText:      random.NewRandomText(rand, 6),
		intIncrementals: make(map[string]*atomic.Int64),
		templates:       make(map[string]*template.Template),
	}
	d.LoadDefaultNgrams()
	return d, nil
}

func NewDecoyWithSeed(seed uint64) (*Decoy, error) {
	if seed == 0 {
		seed = uint64(time.Now().UnixNano())
	}
	return NewDecoy(rand.NewPCG(seed, 1))
}

func (d *Decoy) LoadDefaultNgrams() {
	d.randomText.NgramsFromString(random.DefaultRandomTextCorpus)
}

func (d *Decoy) LoadNgramsString(corpus string) {
	d.randomText.NgramsFromString(corpus)
}

func (d *Decoy) RandomText(maxWords int) (string, error) {
	return d.randomText.RandomText(maxWords)
}

func (d *Decoy) RandomInt(min, max int) int {
	return d.Rand.IntN(max-min) + min
}

func (d *Decoy) RandomFloat(min, max float64) float64 {
	return min + d.Rand.Float64()*(max-min)
}

func (d *Decoy) RandomBoolean() bool {
	x := d.RandomInt(0, 100)
	return x%2 == 0
}

func (d *Decoy) Probability(probability float64) bool {
	x := d.RandomFloat(0, 1)
	return x <= probability
}

func (d *Decoy) CurrentIncrementalInt(id string) (int, error) {
	d.m.RLock()
	defer d.m.RUnlock()
	incremental, exists := d.intIncrementals[id]
	if !exists {
		return 0, fmt.Errorf("Incremental with id %q doesn't exists", id)
	}
	return int(incremental.Load()), nil
}

func (d *Decoy) NextIncrementalInt(id string, start, step int64) int64 {
	d.m.Lock()
	defer d.m.Unlock()
	incremental, exists := d.intIncrementals[id]
	if !exists {
		incremental = &atomic.Int64{}
		incremental.Store(start)
		d.intIncrementals[id] = incremental
	} else {
		d.intIncrementals[id].Add(step)
	}
	return incremental.Load()
}

func (d *Decoy) RandomChoiceAny(choices ...any) any {
	return RandomChoice(d, choices...)
}

func (d *Decoy) ReadFileBytes(path string) ([]byte, error) {
	return os.ReadFile(path)
}

func (d *Decoy) ReadFileBase64(path string) (string, error) {
	bytes, err := d.ReadFileBytes(path)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(bytes), nil
}

func (d *Decoy) ReadFileString(path string) (string, error) {
	bytes, err := os.ReadFile(path)
	return string(bytes), err
}

func (d *Decoy) JsonUnmarshalBytes(data []byte) (map[string]any, error) {
	var res map[string]any
	err := json.Unmarshal(data, &res)
	return res, err
}

func (d *Decoy) JsonUnmarshalString(data string) (map[string]any, error) {
	return d.JsonUnmarshalBytes([]byte(data))
}

func (d *Decoy) list(elems ...any) []any                { return elems }
func (d *Decoy) listString(elems ...string) []string    { return elems }
func (d *Decoy) listInt(elems ...int) []int             { return elems }
func (d *Decoy) listFloat64(elems ...float64) []float64 { return elems }
func (d *Decoy) listBool(elems ...bool) []bool          { return elems }

func coalesce[T comparable](elems ...T) T {
	var zero T
	for _, elem := range elems {
		if elem != zero {
			return elem
		}
	}
	return zero
}

func (d *Decoy) coalesce(elems ...any) any {
	return coalesce(elems...)
}

func (d *Decoy) coalesceString(elems ...string) string {
	return coalesce(elems...)
}

func (d *Decoy) coalesceInt(elems ...int) int {
	return coalesce(elems...)
}

func (d *Decoy) coalesceFloat64(elems ...float64) float64 {
	return coalesce(elems...)
}

func (d *Decoy) newError(msg string, args ...any) (any, error) {
	return nil, fmt.Errorf(msg, args...)
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
		"RandomInt":             d.RandomInt,
		"RandomFloat":           d.RandomFloat,
		"RandomBoolean":         d.RandomBoolean,
		"RandomChoice":          d.RandomChoiceAny,
		"RandomText":            d.RandomText,
		"Probability":           d.Probability,
		"NextIncrementalInt":    d.NextIncrementalInt,
		"CurrentIncrementalInt": d.CurrentIncrementalInt,
		"EnvVariable":           os.Getenv,
		"ReadFileString":        d.ReadFileString,
		"ReadFileBytes":         d.ReadFileBytes,
		"ReadFileBase64":        d.ReadFileBase64,
		"JsonUnmarshalString":   d.JsonUnmarshalString,
		"JsonUnmarshalBytes":    d.JsonUnmarshalBytes,
		"List":                  d.list,
		"ListString":            d.listString,
		"ListInt":               d.listInt,
		"ListFloat64":           d.listFloat64,
		"ListBool":              d.listBool,
		"Coalesce":              d.coalesce,
		"CoalesceString":        d.coalesceString,
		"CoalesceInt":           d.coalesceInt,
		"CoalesceFloat64":       d.coalesceFloat64,
		"NewError":              d.newError,
	}
}

func (d *Decoy) compileTemplate(tmpl string, config *parseTemplateConfig) (*template.Template, error) {
	t, err := template.New(config.Name).
		Funcs(config.Funcs).
		Parse(tmpl)
	if err != nil {
		return nil, err
	}

	if validations.StrIsBlank(config.Name) {
		return t, nil
	}

	d.m.Lock()
	defer d.m.Unlock()

	d.templates[config.Name] = t
	return t, nil
}

func (d *Decoy) CompileTemplate(tmpl string, options ...parseTemplateOption) (*template.Template, error) {
	config := &parseTemplateConfig{
		Funcs: d.DefaultTemplateFuncMaps(),
	}

	for _, option := range options {
		if err := option(config); err != nil {
			return nil, err
		}
	}

	return d.compileTemplate(tmpl, config)
}

func (d *Decoy) ParseTemplate(w io.Writer, tmpl string, options ...parseTemplateOption) error {
	var (
		err      error
		template *template.Template
	)
	config := &parseTemplateConfig{
		Funcs: d.DefaultTemplateFuncMaps(),
	}

	for _, option := range options {
		if err := option(config); err != nil {
			return err
		}
	}

	d.m.RLock()
	template = d.templates[config.Name]
	d.m.RUnlock()
	if template == nil {
		template, err = d.compileTemplate(tmpl, config)
		if err != nil {
			return err
		}
	}

	return template.Execute(w, config.Data)
}

func (d *Decoy) ParseTemplateString(tmpl string, options ...parseTemplateOption) (string, error) {
	str := bytes.NewBufferString("")
	err := d.ParseTemplate(str, tmpl, options...)
	if err != nil {
		return "", err
	}
	return str.String(), nil
}
