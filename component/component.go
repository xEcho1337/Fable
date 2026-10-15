// Package component models Minecraft text components and encodes them as JSON.
package component

import "encoding/json"

type componentKind uint8

const (
	textKind componentKind = iota + 1
	translatableKind
	scoreKind
	selectorKind
	keybindKind
	nbtKind
)

// Component is an immutable Minecraft text component.
type Component struct {
	kind          componentKind
	value         string
	translate     string
	fallback      string
	with          []Component
	score         *ScoreData
	children      []Component
	color         Color
	bold          *bool
	italic        *bool
	underlined    *bool
	strikethrough *bool
	obfuscated    *bool
	clickEvent    *ClickEvent
	hoverEvent    *HoverEvent
	font          string
	insertion     string
	selectorSep   *Component
	nbtSource     string
	nbtInterpret  *bool
	nbtTarget     string
}

// Text creates a literal text component.
func Text(value string) Component {
	return Component{kind: textKind, value: value}
}

// Translatable creates a component resolved by the client's translation key.
func Translatable(key string, args ...Component) Component {
	return Component{kind: translatableKind, translate: key, with: cloneComponents(args)}
}

// Score creates a component displaying a scoreboard value.
func Score(name, objective string) Component {
	return Component{kind: scoreKind, score: &ScoreData{Name: name, Objective: objective}}
}

// ScoreValue creates a scoreboard component with an explicit value.
func ScoreValue(name, objective, value string) Component {
	return Component{kind: scoreKind, score: &ScoreData{Name: name, Objective: objective, Value: value}}
}

// NBT creates a component displaying an NBT path from a source.
func NBT(path, source string, interpret bool) Component {
	return EntityNBT(path, source, interpret)
}

// EntityNBT creates an NBT component reading from an entity selector.
func EntityNBT(path, selector string, interpret bool) Component {
	return Component{kind: nbtKind, value: path, nbtSource: selector, nbtInterpret: boolPointer(interpret), nbtTarget: "entity"}
}

// BlockNBT creates an NBT component reading from a block position.
func BlockNBT(path, position string, interpret bool) Component {
	return Component{kind: nbtKind, value: path, nbtSource: position, nbtInterpret: boolPointer(interpret), nbtTarget: "block"}
}

// StorageNBT creates an NBT component reading from a storage resource.
func StorageNBT(path, storage string, interpret bool) Component {
	return Component{kind: nbtKind, value: path, nbtSource: storage, nbtInterpret: boolPointer(interpret), nbtTarget: "storage"}
}

// Selector creates a component displaying entities matched by a selector.
func Selector(selector string) Component {
	return Component{kind: selectorKind, value: selector}
}

// Keybind creates a component displaying the key bound to the supplied key name.
func Keybind(key string) Component {
	return Component{kind: keybindKind, value: key}
}

// Color sets the component color.
func (c Component) Color(color Color) Component {
	c.color = color
	return c
}

// Fallback sets text used when a translation key is unavailable.
func (c Component) Fallback(value string) Component {
	c.fallback = value
	return c
}

// Bold applies bold decoration.
func (c Component) Bold() Component {
	c.bold = boolPointer(true)
	return c
}

// Italic applies italic decoration.
func (c Component) Italic() Component {
	c.italic = boolPointer(true)
	return c
}

// Underlined applies underlined decoration.
func (c Component) Underlined() Component {
	c.underlined = boolPointer(true)
	return c
}

// Strikethrough applies strikethrough decoration.
func (c Component) Strikethrough() Component {
	c.strikethrough = boolPointer(true)
	return c
}

// Obfuscated applies obfuscated decoration.
func (c Component) Obfuscated() Component {
	c.obfuscated = boolPointer(true)
	return c
}

// SetBold explicitly enables or disables bold decoration.
func (c Component) SetBold(value bool) Component { c.bold = boolPointer(value); return c }

// SetItalic explicitly enables or disables italic decoration.
func (c Component) SetItalic(value bool) Component { c.italic = boolPointer(value); return c }

// SetUnderlined explicitly enables or disables underlined decoration.
func (c Component) SetUnderlined(value bool) Component { c.underlined = boolPointer(value); return c }

// SetStrikethrough explicitly enables or disables strikethrough decoration.
func (c Component) SetStrikethrough(value bool) Component {
	c.strikethrough = boolPointer(value)
	return c
}

// SetObfuscated explicitly enables or disables obfuscated decoration.
func (c Component) SetObfuscated(value bool) Component { c.obfuscated = boolPointer(value); return c }

// Append adds child components after the component's own content.
func (c Component) Append(children ...Component) Component {
	c.children = append(cloneComponents(c.children), children...)
	return c
}

// Empty reports whether the component has no content or style.
func (c Component) Empty() bool {
	return c.kind == 0 && c.value == "" && len(c.children) == 0
}

// Click attaches a click event to the component.
func (c Component) Click(action ClickAction, value string) Component {
	c.clickEvent = &ClickEvent{Action: action, Value: value}
	return c
}

// Hover attaches a hover event to the component.
func (c Component) Hover(action HoverAction, contents any) Component {
	c.hoverEvent = &HoverEvent{Action: action, Contents: contents}
	return c
}

// Font sets the resource location used to render the component.
func (c Component) Font(font string) Component {
	c.font = font
	return c
}

// Insertion sets the text inserted when the component is shift-clicked.
func (c Component) Insertion(value string) Component {
	c.insertion = value
	return c
}

// Separator sets the separator used between selector results.
func (c Component) Separator(separator Component) Component {
	c.selectorSep = &separator
	return c
}

// MarshalJSON encodes the component using Minecraft's JSON component format.
func (c Component) MarshalJSON() ([]byte, error) {
	wire := wireComponent{
		Color:         c.color,
		Bold:          c.bold,
		Italic:        c.italic,
		Underlined:    c.underlined,
		Strikethrough: c.strikethrough,
		Obfuscated:    c.obfuscated,
		Extra:         cloneComponents(c.children),
		ClickEvent:    c.clickEvent,
		HoverEvent:    c.hoverEvent,
		Font:          c.font,
		Insertion:     c.insertion,
		Fallback:      c.fallback,
		Separator:     c.selectorSep,
		NBTInterpret:  c.nbtInterpret,
	}

	switch c.kind {
	case textKind:
		wire.Text = &c.value
	case translatableKind:
		wire.Translate = &c.translate
		wire.With = cloneComponents(c.with)
	case scoreKind:
		wire.Score = c.score
	case selectorKind:
		wire.Selector = &c.value
	case keybindKind:
		wire.Keybind = &c.value
	case nbtKind:
		wire.NBT = &c.value
		switch c.nbtTarget {
		case "block":
			wire.Block = c.nbtSource
		case "storage":
			wire.Storage = c.nbtSource
		default:
			wire.Entity = c.nbtSource
		}
	}

	return json.Marshal(wire)
}

type wireComponent struct {
	Text          *string     `json:"text,omitempty"`
	Translate     *string     `json:"translate,omitempty"`
	Fallback      string      `json:"fallback,omitempty"`
	With          []Component `json:"with,omitempty"`
	Score         *ScoreData  `json:"score,omitempty"`
	Selector      *string     `json:"selector,omitempty"`
	Keybind       *string     `json:"keybind,omitempty"`
	Extra         []Component `json:"extra,omitempty"`
	Color         Color       `json:"color,omitempty"`
	Bold          *bool       `json:"bold,omitempty"`
	Italic        *bool       `json:"italic,omitempty"`
	Underlined    *bool       `json:"underlined,omitempty"`
	Strikethrough *bool       `json:"strikethrough,omitempty"`
	Obfuscated    *bool       `json:"obfuscated,omitempty"`
	ClickEvent    *ClickEvent `json:"clickEvent,omitempty"`
	HoverEvent    *HoverEvent `json:"hoverEvent,omitempty"`
	Font          string      `json:"font,omitempty"`
	Insertion     string      `json:"insertion,omitempty"`
	Separator     *Component  `json:"separator,omitempty"`
	NBT           *string     `json:"nbt,omitempty"`
	Entity        string      `json:"entity,omitempty"`
	Block         string      `json:"block,omitempty"`
	Storage       string      `json:"storage,omitempty"`
	NBTInterpret  *bool       `json:"interpret,omitempty"`
}

func cloneComponents(components []Component) []Component {
	if len(components) == 0 {
		return nil
	}

	return append([]Component(nil), components...)
}

func boolPointer(value bool) *bool {
	return &value
}
