package ch

import (
	"fmt"
)

var (
	// ErrGraphIsFrozen Graph is frozen, so it can not be modified.
	ErrGraphIsFrozen = fmt.Errorf("Graph has been frozen")
)
