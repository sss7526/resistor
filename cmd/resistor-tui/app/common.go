package app

// BaseView provides shared width/height tracking
// for views that depend on terminal dimensions.
//
// Views embed BaseView to avoid repeating
// dimension storage boilerplate.
//
// Resize must be implemented by the view
// if it needs to adjust internal components.
type BaseView struct {
	width  int
	height int
}

// Resize stores the current terminal dimensions.
//
// Views embedding BaseView may extend this method
// to adjust internal Bubble components (e.g., list.SetSize).
func (b *BaseView) Resize(width, height int) {
	b.width = width
	b.height = height
}
