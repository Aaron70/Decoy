# Template Functions

The following functions are available for use inside Go text/templates through the `Decoy` instance. They can be called in any template parsed via `ParseTemplate`, `ParseTemplateString`, or `CompileTemplate`.

> **Sprig functions:** All [Sprig v3](http://masterminds.github.io/sprig/) template functions (date formatting, string manipulation, math, crypto, type conversion, etc.) are also available alongside the functions listed below. For example, Sprig provides `list`, `coalesce`, `env`, `fromJson`, and many other utilities functions.

---

## Random Generation

### `randomInt`
**Signature:** `randomInt(min, max int) (int, error)`  
Returns a random integer in the range `[min, max)`. Returns an error if `min >= max`.

```
{{ randomInt 1 100 }}
{{ randomInt 0 10 }}
```

### `randomFloat`
**Signature:** `randomFloat(min, max float64) (float64, error)`  
Returns a random float64 in the range `[min, max)`. Returns an error if `min > max`.

```
{{ randomFloat 0.0 1.0 }}
```

### `randomBoolean`
**Signature:** `randomBoolean() bool`  
Returns a random boolean value (true/false).

```
{{ randomBoolean }}
```

### `randomChoice`
**Signature:** `randomChoice(args ...any) (any, error)`  
Selects and returns one item at random from the provided list of arguments. Returns an error if no arguments are provided.

```
{{ randomChoice "apple" "banana" "cherry" }}
```

### `randomChoiceList`
**Signature:** `randomChoiceList(choices []any) (any, error)`  
Same as `randomChoice` but takes a slice instead of variadic arguments. Useful when selecting from a list variable built with Sprig's `list` or other sources.

```
{{ $colors := list "red" "green" "blue" }}{{ randomChoiceList $colors }}
```

### `randomText`
**Signature:** `randomText(maxWords int) (string, error)`  
Generates random text using n-gram Markov chains, up to `maxWords` words. Requires n-grams to be loaded (seeded by default corpus).

```
{{ randomText 50 }}
```

### `randomName`
**Signature:** `randomName() string`  
Returns a random first name from a predifined dataset.

```
{{ randomName }}
```

### `randomLastName`
**Signature:** `randomLastName() string`  
Returns a random last name/surname from a predifined dataset.

```
{{ randomLastName }}
```

### `randomFullName`
**Signature:** `randomFullName(middleNameProbability float64) string`  
Combines `randomName` and `randomLastName` into a full name. `middleNameProbability` (0.0 to 1.0) controls the chance of including an additional random name as a middle name.

```
{{ randomFullName 0.0 }}
{{ randomFullName 0.5 }}
{{ randomFullName 1.0 }}
```

---

## Probability

### `probability`
**Signature:** `probability(p float64) bool`  
Returns `true` with probability `p` (0.0 to 1.0), `false` otherwise. Values outside `[0, 1]` are clamped to the nearest bound.

```
{{ probability 0.75 }}
```

---

## Incremental Counters

### `nextIncrementalInt`
**Signature:** `nextIncrementalInt(id string, start, step int64) int64`  
Returns the next value of a named incremental counter. On first call with a given `id`, returns `start` (the `step` parameter is only used on subsequent calls). Each subsequent call adds `step` to the current value.

```
{{ nextIncrementalInt "counter" 1 1 }}
```

> **Note:** The `start` value is returned on the first call; `step` is applied on every subsequent call. Changing `step` between calls will change the increment for that call only.

### `currentIncrementalInt`
**Signature:** `currentIncrementalInt(id string, defaultVal int64) int64`  
Returns the current value of a named incremental counter without advancing it. If the `id` has not been registered yet, returns `defaultVal` instead.

```
{{ currentIncrementalInt "counter" 0 }}
```

---

## I/O Functions

> **Security note:** These functions read arbitrary files from disk at template render time. Only use with trusted template content, as paths are not restricted.

### `readFileString`
**Signature:** `readFileString(path string) (string, error)`  
Reads a file from disk and returns its content as a string.

```
{{ readFileString "data.txt" }}
```

### `readFileBytes`
**Signature:** `readFileBytes(path string) ([]byte, error)`  
Reads a file from disk and returns its raw bytes.

```
{{ readFileBytes "image.png" }}
```

### `readFileBase64`
**Signature:** `readFileBase64(path string) (string, error)`  
Reads a file from disk and returns its content base64-encoded.

```
{{ readFileBase64 "image.png" }}
```
