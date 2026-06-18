package parser

import "github.com/srdino/dino-hatch/internal/ast"

// inheritProps hereda propiedades del padre si el hijo no las definio.
func inheritProps(s *ast.ComputedStyle, parent *ast.ComputedStyle, exp map[string]bool) {
	if !exp["color"] {
		s.Color = parent.Color
	}
	if !exp["bg"] {
		s.BgColor = parent.BgColor
	}
	if !exp["direction"] {
		s.Direction = parent.Direction
	}
	if !exp["align"] {
		s.Align = parent.Align
	}
	if !exp["justify"] {
		s.Justify = parent.Justify
	}
	if !exp["text-align"] {
		s.TextAlign = parent.TextAlign
	}
	if !exp["valign"] {
		s.VAlign = parent.VAlign
	}
	if !exp["overflow"] {
		s.Overflow = parent.Overflow
	}
}
