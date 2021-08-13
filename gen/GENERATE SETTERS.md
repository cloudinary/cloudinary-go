# Typed setters code generation for structs

### Example:

Assuming we have a struct `TestType` with a couple of fields that could store a value with different types.

```go
type TestType struct {
	field1       interface{}
	field2      interface{}
}
```

We want to limit available types with typed setters and add some fluent-style interface for it.

**The final interface should be like this:**
```go
t := TestType{}
t.Field1(1).Field2Percent(0.85)

t1 := TestType{}
t1.Field2Expr("Test")
```

`field1` should be int, `field2` should be `float32` (when we want to set it to percent) or `string` (if we want to pass an expression here).
For our interface to be clear it'd be great to have specific suffixes for named values (like percent or expression).

**We could create typed setters by ourselves:**
```go
func (t *TestType) Field1(field1 int) *TestType {
	t.field1 = field1
	
	return t
}

func (t *TestType) Field2Percent(field2 float32) *TestType {
t.field2 = field2

return t
}

func (t *TestType) Field2Expr(field2 string) *TestType {
t.field2 = field2

return t
}
```

But what if we need this type of setters for multiple structs?
Golang provides type embedding for these needs, but this approach does not work with fluent interfaces as the return type would be an embed struct, not a high-level struct.

We need to create these setters for every struct.

> Let's use the code generator for it!

**Changes required to the struct:**
```go
type TestType struct {
	field1       interface{} `setters:"int"`
	field2      interface{} `setters:"float32:Percent,string:Expr"`
}
```

Run `go run gen/generate_setters.go` then. `testtype_setters.go` would be generated and put near the original file.

Now `TestType` has all needed setters with suffixes.

### Annotation structure

`setters:"type:suffix"`
Types should be divided with comma.
When suffix is provided it would be used in the function name.
If suffix is not set function name would contain only parameter name without any suffix.

### Code generator

`gen/generate_setters.go` walks through all `*.go` files and searches for `setters` annotations in the structs.
- All structs are placed in Directed Acyclic Graph to calculate its root elements.
- Setters are generated only for root elements.
- Setters would be generated even for cross-module embedding
- When cross-module embedding only exported fields should be used for setters

### Example

```go
type TestType struct {
field1       interface{} `setters:"int"`
field2      interface{} `setters:"float32:Percent,string:Expr"`
}

type RootType struct {
	TestType
}
```

In this case setters file would be generated only for `RootType`.