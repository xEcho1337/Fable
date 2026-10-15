package component

import "strings"

// Color is a Minecraft named color or a custom hex color such as "#55ffff".
type Color string

const (
	Black       Color = "black"
	DarkBlue    Color = "dark_blue"
	DarkGreen   Color = "dark_green"
	DarkAqua    Color = "dark_aqua"
	DarkRed     Color = "dark_red"
	DarkPurple  Color = "dark_purple"
	Gold        Color = "gold"
	Gray        Color = "gray"
	DarkGray    Color = "dark_gray"
	Blue        Color = "blue"
	Green       Color = "green"
	Aqua        Color = "aqua"
	Red         Color = "red"
	LightPurple Color = "light_purple"
	Yellow      Color = "yellow"
	White       Color = "white"
)

// CustomColor creates a color from a Minecraft-compatible value.
func CustomColor(value string) Color {
	return Color(value)
}

// IsEmpty reports whether no color is set.
func (c Color) IsEmpty() bool {
	return strings.TrimSpace(string(c)) == ""
}

// ScoreData identifies a scoreboard value.
type ScoreData struct {
	Name      string `json:"name"`
	Objective string `json:"objective"`
	Value     string `json:"value,omitempty"`
}

// ItemHoverContents describes an item displayed by a hover event.
type ItemHoverContents struct {
	ID    string `json:"id"`
	Count int    `json:"count,omitempty"`
}

// EntityHoverContents describes an entity displayed by a hover event.
type EntityHoverContents struct {
	ID   string    `json:"id"`
	Name Component `json:"name,omitempty"`
	Type string    `json:"type,omitempty"`
}
