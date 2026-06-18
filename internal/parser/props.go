package parser

import (
	"log"
	"strings"

	"github.com/srdino/dino-hatch/internal/ast"
)

// resolveCSSVars reemplaza ocurrencias de var(--nombre) con el valor del mapa.
// Si la variable no existe, deja var(--nombre) sin reemplazar (el parser lo ignorara).
func resolveCSSVars(val string, vars map[string]string) string {
	if vars == nil || !strings.Contains(val, "var(--") {
		return val
	}
	for k, v := range vars {
		placeholder := "var(" + k + ")"
		val = strings.ReplaceAll(val, placeholder, v)
	}
	return val
}

// applyProps aplica un mapa de propiedades a un ComputedStyle.
// Marca en `exp` las propiedades aplicadas para controlar herencia.
// `vars` son las variables HSS de :root para resolver var(--name).
func applyProps(s *ast.ComputedStyle, props map[string]string, vars map[string]string, exp map[string]bool) {
	for k, v := range props {
		v = resolveCSSVars(v, vars)
		switch k {
		case "width":
			if n, ok := parseInt(v); ok {
				s.Width = n
				exp["width"] = true
			}
		case "height":
			if n, ok := parseInt(v); ok {
				s.Height = n
				exp["height"] = true
			}
		case "grow":
			if f, ok := parseFloat(v); ok {
				s.Grow = f
				exp["grow"] = true
			} else {
				log.Printf("compute: warning: grow invalido '%s'", v)
			}
		case "gap":
			if n, ok := parseInt(v); ok {
				s.Gap = n
				exp["gap"] = true
			}
		case "direction":
			s.Direction = v
			exp["direction"] = true
		case "align":
			s.Align = v
			exp["align"] = true
		case "justify":
			s.Justify = v
			exp["justify"] = true
		case "color":
			if c, ok := parseColor(v); ok {
				s.Color = c
				exp["color"] = true
			}
		case "bg":
			if c, ok := parseColor(v); ok {
				s.BgColor = c
				exp["bg"] = true
			}
		case "border":
			if b, ok := parseBool(v); ok {
				s.Border = b
				exp["border"] = true
			}
		case "padding":
			if n, ok := parseInt(v); ok {
				s.Padding = n
				exp["padding"] = true
			}
		case "margin":
			if n, ok := parseInt(v); ok {
				s.Margin = n
				exp["margin"] = true
			}
		case "min-width":
			if n, ok := parseInt(v); ok {
				s.MinWidth = n
				exp["min-width"] = true
			}
		case "min-height":
			if n, ok := parseInt(v); ok {
				s.MinHeight = n
				exp["min-height"] = true
			}
		case "max-width":
			if n, ok := parseInt(v); ok {
				s.MaxWidth = n
				exp["max-width"] = true
			}
		case "max-height":
			if n, ok := parseInt(v); ok {
				s.MaxHeight = n
				exp["max-height"] = true
			}
		case "text-align":
			s.TextAlign = v
			exp["text-align"] = true
		case "valign":
			s.VAlign = v
			exp["valign"] = true
		case "overflow":
			s.Overflow = v
			exp["overflow"] = true
		}
	}
}


