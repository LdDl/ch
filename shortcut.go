package ch

// ShortcutPath Representation of shortcut between vertices
//
// From - ID of source vertex
// To - ID of target vertex
// ViaVertex - ID of vertex through which the shortcut exists
// Cost - summary cost of path between two vertices
//
type ShortcutPath struct {
	From int64
	To   int64
	Via  int64
	Cost float64
}
