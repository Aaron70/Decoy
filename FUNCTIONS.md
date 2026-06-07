# Template Functions

The following functions are available for use inside Go text/templates through the `Decoy` instance. They can be called in any template parsed via `ParseTemplate`, `ParseTemplateString`, or `CompileTemplate`.

---

## Random Generation

### `RandomInt`
**Signature:** `RandomInt(min, max int) (int, error)`  
Returns a random integer in the range `[min, max)`. Returns an error if `min >= max`.

```
{{ RandomInt 1 100 }}
{{ RandomInt 0 10 }}
```

### `RandomFloat`
**Signature:** `RandomFloat(min, max float64) float64`  
Returns a random float64 in the range `[min, max)`.

```
{{ RandomFloat 0.0 1.0 }}
```

### `RandomBoolean`
**Signature:** `RandomBoolean() bool`  
Returns a random boolean value (true/false).

```
{{ RandomBoolean }}
```

### `RandomChoice`
**Signature:** `RandomChoice(args ...any) (any, error)`  
Selects and returns one item at random from the provided list of arguments. Returns an error if no arguments are provided.

```
{{ RandomChoice "apple" "banana" "cherry" }}
```

### `RandomText`
**Signature:** `RandomText(maxWords int) (string, error)`  
Generates random text using n-gram Markov chains, up to `maxWords` words. Requires n-grams to be loaded (seeded by default corpus).

```
{{ RandomText 50 }}
```

---

## Probability

### `Probability`
**Signature:** `Probability(probability float64) bool`  
Returns `true` with probability `p` (0.0 to 1.0), `false` otherwise. Internally uses `x < p` where `x` is a random float in `[0.0, 1.0)`, so `Probability(0.0)` never returns true.

```
{{ Probability 0.75 }}
```

---

## Incremental Counters

### `NextIncrementalInt`
**Signature:** `NextIncrementalInt(id string, start, step int64) int64`  
Returns the next value of a named incremental counter. On first call with a given `id`, returns `start` (the `step` parameter is only used on subsequent calls). Each subsequent call adds `step` to the current value.

```
{{ NextIncrementalInt "counter" 1 1 }}
```

> **Note:** The `start` value is returned on the first call; `step` is applied on every subsequent call. Changing `step` between calls will change the increment for that call only.

### `CurrentIncrementalInt`
**Signature:** `CurrentIncrementalInt(id string, defaultVal int) int`  
Returns the current value of a named incremental counter without advancing it. If the `id` has not been registered yet, returns `defaultVal` instead.

```
{{ CurrentIncrementalInt "counter" 0 }}
```

---

## Environment

### `EnvVariable`
**Signature:** `EnvVariable(key string) string`  
Reads an environment variable by name. Returns an empty string if not set.

```
{{ EnvVariable "HOME" }}
```

---

## I/O Functions

> **Security note:** These functions read arbitrary files from disk at template render time. Only use with trusted template content, as paths are not restricted.

### `ReadFileString`
**Signature:** `ReadFileString(path string) (string, error)`  
Reads a file from disk and returns its content as a string.

```
{{ ReadFileString "data.txt" }}
```

### `ReadFileBytes`
**Signature:** `ReadFileBytes(path string) ([]byte, error)`  
Reads a file from disk and returns its raw bytes.

```
{{ ReadFileBytes "image.png" }}
```

### `ReadFileBase64`
**Signature:** `ReadFileBase64(path string) (string, error)`  
Reads a file from disk and returns its content base64-encoded.

```
{{ ReadFileBase64 "image.png" }}
```

---

## Serialization/Deserialization

### `JsonUnmarshalString`
**Signature:** `JsonUnmarshalString(data string) (map[string]any, error)`  
Parses a JSON string into a map for further template access.

```
{{ (JsonUnmarshalString `{"name": "Alice"}`).name }}
```

### `JsonUnmarshalBytes`
**Signature:** `JsonUnmarshalBytes(data []byte) (map[string]any, error)`  
Parses JSON bytes into a map.

```
{{ $data := ReadFileBytes "config.json" }}{{ (JsonUnmarshalBytes $data).key }}
```

---

## List Builders

### `List`
**Signature:** `List(args ...any) []any`  
Creates a slice of mixed-type values from the given arguments.

```
{{ range List "a" 1 true }}{{ . }} {{ end }}
```

### `ListString`
**Signature:** `ListString(args ...string) []string`  
Creates a string slice.

```
{{ range ListString "x" "y" "z" }}{{ . }} {{ end }}
```

### `ListInt`
**Signature:** `ListInt(args ...int) []int`  
Creates an int slice.

```
{{ range ListInt 10 20 30 }}{{ . }} {{ end }}
```

### `ListFloat64`
**Signature:** `ListFloat64(args ...float64) []float64`  
Creates a float64 slice.

```
{{ range ListFloat64 1.5 2.5 3.5 }}{{ . }} {{ end }}
```

### `ListBool`
**Signature:** `ListBool(args ...bool) []bool`  
Creates a bool slice.

```
{{ range ListBool true false true }}{{ . }} {{ end }}
```

---

## Coalesce (first non-zero)

### `Coalesce`
**Signature:** `Coalesce(args ...any) any`  
Returns the first non-zero argument. Useful for fallback/default values.

```
{{ Coalesce .Title "Default Title" }}
```

### `CoalesceString`
**Signature:** `CoalesceString(args ...string) string`  
String-typed variant of Coalesce.

```
{{ CoalesceString .Name "unknown" }}
```

### `CoalesceInt`
**Signature:** `CoalesceInt(args ...int) int`  
Int-typed variant of Coalesce.

```
{{ CoalesceInt .Count 0 }}
```

### `CoalesceFloat64`
**Signature:** `CoalesceFloat64(args ...float64) float64`  
Float64-typed variant of Coalesce.

```
{{ CoalesceFloat64 .Price 9.99 }}
```

---

## Error

### `NewError`
**Signature:** `NewError(msg string, args ...any) (any, error)`  
Returns an error with a formatted message, halting template execution. Useful for validation or pre-condition checks.

```
{{ NewError "unexpected value: %v" .Value }}
```
