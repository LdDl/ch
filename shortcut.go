package ch

// ShortcutInfo Information about shortcut
//
// ViaVertex - ID of vertex through which the contraction exists
// Cost - summary cost of path between two vertices
//
type ShortcutInfo struct {
	ViaVertex int
	Cost      float64
}
