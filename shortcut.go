package ch

// ShortcutPath Representation of shortcut between vertices
type ShortcutPath struct {
	From int64
	To   int64
	Via  int64
	Cost float64
}
