package tcutil

// Color is the color of the component.
//
//tgcomp:export
type Color string

//tgcomp:export
const (
	// ColorNull is the default color.
	ColorNull Color = ""

	ColorInfo    Color = "info"
	ColorSuccess Color = "success"
	ColorWarning Color = "warning"
	ColorDanger  Color = "danger"
)
