package ast

import "github.com/gdamore/tcell/v2"

// Document representa el archivo .hml completo
type Document struct {
	Pages     []Page
	ThemeVars map[string]string // CSS variables de :root, ej: {"--bg": "#333"}
}

// Page representa una pantalla/interfaz
type Page struct {
	Name     string
	Width    int          // 0 = auto/terminal width
	Height   int          // 0 = auto/terminal height
	Style    ComputedStyle
	Children []ElementNode
}

// ElementNode representa un elemento UI
type ElementNode struct {
	Tag        string            // "box", "text", etc.
	Attrs      map[string]string // atributos raw del XML
	Style      ComputedStyle
	Children   []ElementNode
	Text       string              // contenido textual para <text>
	BoundBox   BoundBox            // calculado en fase layout
	ID         string              // atributo id
	ScrollX    int                 // scroll offset horizontal
	ScrollY    int                 // scroll offset vertical
	Events     map[string]string   // "click" -> "page:proyectos"
	IncludeSrc string              // "" si no es include, "<path>" si es include
}

// ComputedStyle con todas las propiedades resueltas
type ComputedStyle struct {
	Width     int
	Height    int
	Grow      float64
	Gap       int
	Align     string // "start", "center", "end", "stretch"
	Justify   string // "start", "center", "end", "space-between"
	Direction string // "row", "column"
	Color     tcell.Color
	BgColor   tcell.Color
	Border    bool
	TextAlign string    // "left", "center", "right"
	VAlign    string    // "top", "middle", "bottom"
	Padding   int       // espacio entre borde y contenido interno
	Margin    int       // espacio externo entre hermanos
	MinWidth  int       // 0 = sin constraint
	MinHeight int       // 0 = sin constraint
	MaxWidth  int       // 0 = sin constraint
	MaxHeight int       // 0 = sin constraint
	Overflow  string    // "visible" | "hidden" | "scroll" (default "visible")
}

// BoundBox posicion calculada por layout
type BoundBox struct {
	X, Y, W, H int
}

// StyleRule representa una regla HSS sin resolver
type StyleRule struct {
	Selector   string
	Properties map[string]string
}

// DefaultStyle devuelve un ComputedStyle con valores por defecto
func DefaultStyle() ComputedStyle {
	return ComputedStyle{
		Width:     0,
		Height:    0,
		Grow:      0,
		Gap:       0,
		Align:     "stretch",
		Justify:   "start",
		Direction: "column",
		Color:     tcell.ColorWhite,
		BgColor:   tcell.ColorReset,
		Border:    false,
		TextAlign: "center",
		VAlign:    "middle",
		Overflow:  "visible",
	}
}
