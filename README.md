# Fable

Fable is the Relay-independent Minecraft text component library. The public
package is named `component` so plugin code can build messages fluently:

```go
import "github.com/xEcho1337/fable/component"

message := component.Text("Ciao").
	Color(component.Green).
	Bold().
	Append(component.Text("!"))
```

Components implement `json.Marshaler` and produce Minecraft's JSON component
format. The package has no dependency on Relay and can be extracted into its
own module later without changing the component API.

`Audience` is only an interface; the backend decides how a message reaches a
player or console. The package also provides `ParsePlain`, `ParseLegacy`,
`ParseMiniMessage`, `Plain`, `Legacy`, and `MiniMessage` for common text
formats. The MiniMessage implementation intentionally covers the core style
tags and preserves unknown tags as text.

MiniMessage also supports click and hover tags, hex RGB colors (`<#55ffff>`),
and linear gradients (`<gradient:#ff0000:#0000ff>text</gradient>`).

The default encoder uses Go's standard `encoding/json`, which is portable on
all supported targets, including TinyGo/WASM. Applications that need a SIMD
JSON encoder can use that encoder directly with `Component.MarshalJSON`; Fable
keeps the wire model independent of a particular JSON implementation.
