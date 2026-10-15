package component

// ClickAction identifies the operation performed by a click event.
type ClickAction string

const (
	OpenURL         ClickAction = "open_url"
	OpenFile        ClickAction = "open_file"
	RunCommand      ClickAction = "run_command"
	SuggestCommand  ClickAction = "suggest_command"
	ChangePage      ClickAction = "change_page"
	CopyToClipboard ClickAction = "copy_to_clipboard"
)

// HoverAction identifies the content shown by a hover event.
type HoverAction string

const (
	ShowText   HoverAction = "show_text"
	ShowItem   HoverAction = "show_item"
	ShowEntity HoverAction = "show_entity"
)

// ClickEvent describes a client-side click action.
type ClickEvent struct {
	Action ClickAction `json:"action"`
	Value  string      `json:"value"`
}

// HoverEvent describes content shown while the cursor hovers over a component.
type HoverEvent struct {
	Action   HoverAction `json:"action"`
	Contents any         `json:"contents"`
}

// Item creates hover contents for an item stack.
func Item(id string, count int) ItemHoverContents {
	return ItemHoverContents{ID: id, Count: count}
}

// Entity creates hover contents for an entity.
func Entity(id, entityType string, name Component) EntityHoverContents {
	return EntityHoverContents{ID: id, Type: entityType, Name: name}
}
