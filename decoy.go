package decoy

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"math/rand/v2"
	"os"
	"sync"
	"sync/atomic"
	"text/template"
	"time"

	"github.com/Masterminds/sprig/v3"
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

func (d *Decoy) RandomName() string {
	return d.randomText.RandomName()
}

func (d *Decoy) RandomLastName() string {
	return d.randomText.RandomLastName()
}

func (d *Decoy) RandomFullName(middleNameProbability float64) string {
	name := d.RandomName()
	lastName := d.RandomLastName()
	middle := d.Probability(middleNameProbability)
	if middle {
		return fmt.Sprintf("%s %s %s", name, d.RandomLastName(), lastName)
	}

	return fmt.Sprintf("%s %s", name, lastName)
}

func (d *Decoy) RandomInt(min, max int) (int, error) {
	if min >= max {
		return 0, fmt.Errorf("min (%d) must be less than max (%d)", min, max)
	}
	return d.Rand.IntN(max-min) + min, nil
}

func (d *Decoy) RandomFloat(min, max float64) (float64, error) {
	if min > max {
		return 0, fmt.Errorf("RandomFloat: min (%f) must be less than or equal to max (%f)", min, max)
	}
	return min + d.Rand.Float64()*(max-min), nil
}

func (d *Decoy) RandomBoolean() bool {
	return d.Rand.IntN(2) == 0
}

func (d *Decoy) Probability(probability float64) bool {
	if probability <= 0 {
		return false
	}
	if probability >= 1 {
		return true
	}
	x, _ := d.RandomFloat(0, 1)
	return x < probability
}

func (d *Decoy) CurrentIncrementalInt(id string, defaultVal int64) int64 {
	d.m.RLock()
	defer d.m.RUnlock()
	incremental, exists := d.intIncrementals[id]
	if !exists {
		return defaultVal
	}
	return incremental.Load()
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

func (d *Decoy) RandomChoiceAny(choices ...any) (any, error) {
	return RandomChoice(d, choices...)
}

func (d *Decoy) RandomChoiceList(choices []any) (any, error) {
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
	funcMap := template.FuncMap{
		"randomInt":             d.RandomInt,
		"randomFloat":           d.RandomFloat,
		"randomBoolean":         d.RandomBoolean,
		"randomChoice":          d.RandomChoiceAny,
		"randomChoiceList":      d.RandomChoiceList,
		"randomText":            d.RandomText,
		"randomName":            d.RandomName,
		"randomLastName":        d.RandomLastName,
		"randomFullName":        d.RandomFullName,
		"probability":           d.Probability,
		"nextIncrementalInt":    d.NextIncrementalInt,
		"currentIncrementalInt": d.CurrentIncrementalInt,
		"readFileString":        d.ReadFileString,
		"readFileBytes":         d.ReadFileBytes,
		"readFileBase64":        d.ReadFileBase64,
	}
	maps.Insert(funcMap, maps.All(sprig.FuncMap()))
	return funcMap
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
