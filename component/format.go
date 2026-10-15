package component

import (
	"fmt"
	"strconv"
	"strings"
)

// ParsePlain converts plain text into a text component.
func ParsePlain(value string) Component {
	return Text(value)
}

// Plain returns the textual content of a component tree.
func Plain(component Component) string {
	var builder strings.Builder
	writePlain(&builder, component)
	return builder.String()
}

func writePlain(builder *strings.Builder, component Component) {
	switch component.kind {
	case textKind:
		builder.WriteString(component.value)
	case translatableKind:
		if component.fallback != "" {
			builder.WriteString(component.fallback)
		}
	case selectorKind:
		builder.WriteString(component.value)
	case keybindKind:
		builder.WriteString(component.value)
	case scoreKind:
		if component.score != nil {
			builder.WriteString(component.score.Value)
		}
	case nbtKind:
		builder.WriteString(component.value)
	}

	for _, child := range component.children {
		writePlain(builder, child)
	}
}

// ParseLegacy parses Minecraft legacy formatting codes. Both § and & are
// accepted as prefixes; unsupported codes are retained as literal text.
func ParseLegacy(value string) Component {
	return parseLegacy(value, true)
}

func parseLegacy(value string, acceptAmpersand bool) Component {
	value = strings.ReplaceAll(value, "\u00a7", "&")
	var parts []Component
	var text strings.Builder
	style := legacyStyle{}
	flush := func() {
		if text.Len() == 0 {
			return
		}

		parts = append(parts, style.apply(Text(text.String())))
		text.Reset()
	}

	for index := 0; index < len(value); index++ {
		if value[index] != '&' || !acceptAmpersand {
			text.WriteByte(value[index])
			continue
		}

		if index+1 >= len(value) {
			text.WriteByte(value[index])
			continue
		}

		code := strings.ToLower(string(value[index+1]))
		if !style.applyCode(code) {
			text.WriteByte(value[index])
			continue
		}

		flush()
		index++
	}

	flush()
	return join(parts)
}

// Legacy serializes a component tree using § formatting codes.
func Legacy(component Component) string {
	var builder strings.Builder
	writeLegacy(&builder, component, legacyStyle{})
	return builder.String()
}

func writeLegacy(builder *strings.Builder, component Component, inherited legacyStyle) {
	style := inherited.merge(component)
	if style != inherited {
		builder.WriteString("\u00a7r")
		builder.WriteString(style.codes())
	}

	builder.WriteString(Plain(componentWithoutChildren(component)))
	for _, child := range component.children {
		writeLegacy(builder, child, style)
	}
}

// ParseMiniMessage parses the common MiniMessage tags used for text styling.
// Unknown tags remain literal text, preserving information for callers.
func ParseMiniMessage(value string) Component {
	root := Text("")
	stack := []Component{root}
	tags := []string{""}
	text := strings.Builder{}
	flush := func() {
		if text.Len() == 0 {
			return
		}

		parent := stack[len(stack)-1]
		stack[len(stack)-1] = parent.Append(Text(text.String()))
		text.Reset()
	}

	for index := 0; index < len(value); {
		if value[index] != '<' {
			text.WriteByte(value[index])
			index++
			continue
		}

		end := strings.IndexByte(value[index:], '>')
		if end < 0 {
			text.WriteString(value[index:])
			break
		}

		tag := value[index+1 : index+end]
		if !applyMiniTag(&stack, tag) {
			text.WriteString(value[index : index+end+1])
			index += end + 1
			continue
		}

		flush()
		if strings.HasPrefix(tag, "/") {
			child := stack[len(stack)-1]
			stack = stack[:len(stack)-1]
			openingTag := tags[len(tags)-1]
			tags = tags[:len(tags)-1]
			if strings.HasPrefix(openingTag, "gradient:") {
				child = applyGradientTag(child, openingTag)
			}
			parent := stack[len(stack)-1]
			stack[len(stack)-1] = parent.Append(child)
		} else {
			tags = append(tags, tag)
		}
		index += end + 1
	}

	flush()
	for len(stack) > 1 {
		child := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		tags = tags[:len(tags)-1]
		stack[len(stack)-1] = stack[len(stack)-1].Append(child)
	}

	return stack[0]
}

// MiniMessage serializes basic text styles using MiniMessage tags.
func MiniMessage(component Component) string {
	var builder strings.Builder
	writeMiniMessage(&builder, component)
	return builder.String()
}

func writeMiniMessage(builder *strings.Builder, component Component) {
	if component.kind == textKind {
		writeMiniStyledText(builder, component)
	} else {
		builder.WriteString(escapeMini(Plain(componentWithoutChildren(component))))
	}

	for _, child := range component.children {
		writeMiniMessage(builder, child)
	}
}

func writeMiniStyledText(builder *strings.Builder, component Component) {
	value := escapeMini(component.value)
	if component.color != "" {
		value = "<" + string(component.color) + ">" + value + "</" + string(component.color) + ">"
	}
	if component.bold != nil && *component.bold {
		value = "<bold>" + value + "</bold>"
	}
	if component.italic != nil && *component.italic {
		value = "<italic>" + value + "</italic>"
	}
	if component.underlined != nil && *component.underlined {
		value = "<underlined>" + value + "</underlined>"
	}
	if component.strikethrough != nil && *component.strikethrough {
		value = "<strikethrough>" + value + "</strikethrough>"
	}
	if component.obfuscated != nil && *component.obfuscated {
		value = "<obfuscated>" + value + "</obfuscated>"
	}
	builder.WriteString(value)
}

type legacyStyle struct {
	color         Color
	bold          bool
	italic        bool
	underlined    bool
	strikethrough bool
	obfuscated    bool
}

func (style *legacyStyle) applyCode(code string) bool {
	colors := map[string]Color{"0": Black, "1": DarkBlue, "2": DarkGreen, "3": DarkAqua, "4": DarkRed, "5": DarkPurple, "6": Gold, "7": Gray, "8": DarkGray, "9": Blue, "a": Green, "b": Aqua, "c": Red, "d": LightPurple, "e": Yellow, "f": White}
	if color, ok := colors[code]; ok {
		*style = legacyStyle{color: color}
		return true
	}

	switch code {
	case "k":
		style.obfuscated = true
	case "l":
		style.bold = true
	case "m":
		style.strikethrough = true
	case "n":
		style.underlined = true
	case "o":
		style.italic = true
	case "r":
		*style = legacyStyle{}
	default:
		return false
	}
	return true
}

func (style legacyStyle) apply(component Component) Component {
	return component.Color(style.color).SetBold(style.bold).SetItalic(style.italic).SetUnderlined(style.underlined).SetStrikethrough(style.strikethrough).SetObfuscated(style.obfuscated)
}

func (style legacyStyle) merge(component Component) legacyStyle {
	if component.color != "" {
		style.color = component.color
	}
	if component.bold != nil {
		style.bold = *component.bold
	}
	if component.italic != nil {
		style.italic = *component.italic
	}
	if component.underlined != nil {
		style.underlined = *component.underlined
	}
	if component.strikethrough != nil {
		style.strikethrough = *component.strikethrough
	}
	if component.obfuscated != nil {
		style.obfuscated = *component.obfuscated
	}
	return style
}

func (style legacyStyle) codes() string {
	var result strings.Builder
	colors := map[Color]string{Black: "0", DarkBlue: "1", DarkGreen: "2", DarkAqua: "3", DarkRed: "4", DarkPurple: "5", Gold: "6", Gray: "7", DarkGray: "8", Blue: "9", Green: "a", Aqua: "b", Red: "c", LightPurple: "d", Yellow: "e", White: "f"}
	if code := colors[style.color]; code != "" {
		result.WriteString("\u00a7" + code)
	}
	if style.bold {
		result.WriteString("\u00a7l")
	}
	if style.italic {
		result.WriteString("\u00a7o")
	}
	if style.underlined {
		result.WriteString("\u00a7n")
	}
	if style.strikethrough {
		result.WriteString("\u00a7m")
	}
	if style.obfuscated {
		result.WriteString("\u00a7k")
	}
	return result.String()
}

func applyMiniTag(stack *[]Component, tag string) bool {
	if strings.HasPrefix(tag, "/") {
		return len(*stack) > 1
	}
	component := Text("")
	switch {
	case tag == "bold":
		component = component.Bold()
	case tag == "italic":
		component = component.Italic()
	case tag == "underlined":
		component = component.Underlined()
	case tag == "strikethrough":
		component = component.Strikethrough()
	case tag == "obfuscated":
		component = component.Obfuscated()
	case strings.HasPrefix(tag, "color:"):
		component = component.Color(CustomColor(strings.TrimPrefix(tag, "color:")))
	case strings.HasPrefix(tag, "#"):
		component = component.Color(CustomColor(tag))
	case strings.HasPrefix(tag, "font:"):
		component = component.Font(strings.TrimPrefix(tag, "font:"))
	case strings.HasPrefix(tag, "click:"):
		parts := strings.SplitN(strings.TrimPrefix(tag, "click:"), ":", 2)
		if len(parts) != 2 {
			return false
		}
		component = component.Click(ClickAction(parts[0]), strings.Trim(parts[1], "'\""))
	case strings.HasPrefix(tag, "hover:"):
		parts := strings.SplitN(strings.TrimPrefix(tag, "hover:"), ":", 2)
		if len(parts) != 2 {
			return false
		}
		component = component.Hover(HoverAction(parts[0]), ParseMiniMessage(strings.Trim(parts[1], "'\"")))
	case strings.HasPrefix(tag, "gradient:"):
		if len(strings.Split(strings.TrimPrefix(tag, "gradient:"), ":")) != 2 {
			return false
		}
	default:
		colors := map[string]Color{"black": Black, "dark_blue": DarkBlue, "dark_green": DarkGreen, "dark_aqua": DarkAqua, "dark_red": DarkRed, "dark_purple": DarkPurple, "gold": Gold, "gray": Gray, "dark_gray": DarkGray, "blue": Blue, "green": Green, "aqua": Aqua, "red": Red, "light_purple": LightPurple, "yellow": Yellow, "white": White}
		color, ok := colors[tag]
		if !ok {
			return false
		}
		component = component.Color(color)
	}
	*stack = append(*stack, component)
	return true
}

func applyGradientTag(component Component, tag string) Component {
	colors := strings.Split(strings.TrimPrefix(tag, "gradient:"), ":")
	if len(colors) != 2 {
		return component
	}

	start, startOK := parseRGB(colors[0])
	end, endOK := parseRGB(colors[1])
	if !startOK || !endOK {
		return component
	}

	value := Plain(component)
	runes := []rune(value)
	if len(runes) == 0 {
		return component
	}

	parts := make([]Component, 0, len(runes))
	for index, character := range runes {
		factor := float64(index) / float64(max(1, len(runes)-1))
		parts = append(parts, Text(string(character)).Color(CustomColor(formatRGB(interpolateRGB(start, end, factor)))))
	}
	return join(parts)
}

type rgb struct{ red, green, blue uint8 }

func parseRGB(value string) (rgb, bool) {
	value = strings.TrimPrefix(value, "#")
	if len(value) != 6 {
		return rgb{}, false
	}
	red, errRed := strconv.ParseUint(value[0:2], 16, 8)
	green, errGreen := strconv.ParseUint(value[2:4], 16, 8)
	blue, errBlue := strconv.ParseUint(value[4:6], 16, 8)
	if errRed != nil || errGreen != nil || errBlue != nil {
		return rgb{}, false
	}
	return rgb{uint8(red), uint8(green), uint8(blue)}, true
}

func interpolateRGB(start, end rgb, factor float64) rgb {
	interpolate := func(left, right uint8) uint8 {
		return uint8(float64(left) + (float64(right)-float64(left))*factor)
	}
	return rgb{interpolate(start.red, end.red), interpolate(start.green, end.green), interpolate(start.blue, end.blue)}
}

func formatRGB(value rgb) string {
	return fmt.Sprintf("#%02x%02x%02x", value.red, value.green, value.blue)
}

func max(left, right int) int {
	if left > right {
		return left
	}
	return right
}

func join(parts []Component) Component {
	if len(parts) == 0 {
		return Text("")
	}
	return parts[0].Append(parts[1:]...)
}

func componentWithoutChildren(component Component) Component {
	component.children = nil
	return component
}

func escapeMini(value string) string {
	return strings.NewReplacer("\\", "\\\\", "<", "\\<", ">", "\\>").Replace(value)
}
