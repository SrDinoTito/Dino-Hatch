// Package parser implementa el parseo de archivos .hml.
package parser

import (
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

// ComputeStyles recibe un Document (AST raw), las reglas HSS,
// y las CSS variables de :root, y resuelve los estilos computados para cada nodo.
func ComputeStyles(doc *ast.Document, rules []ast.StyleRule, vars map[string]string) *ast.Document {
	rulesBySel := make(map[string]map[string]string, len(rules))
	for _, r := range rules {
		rulesBySel[r.Selector] = r.Properties
	}

	for i := range doc.Pages {
		p := &doc.Pages[i]
		p.Style = ast.DefaultStyle()
		resolveNode(&p.Children, rulesBySel, vars, &p.Style)
	}
	return doc
}

// resolveNode aplica estilos recursivamente: defaults, HSS, inline, herencia.
// `vars` son las CSS variables de :root para resolver var(--name).
func resolveNode(children *[]ast.ElementNode, rulesBySel map[string]map[string]string, vars map[string]string, parent *ast.ComputedStyle) {
	for i := range *children {
		n := &(*children)[i]

		s := ast.DefaultStyle()
		exp := make(map[string]bool)

		if props, ok := rulesBySel[n.Tag]; ok {
			applyProps(&s, props, vars, exp)
		}
		applyProps(&s, n.Attrs, vars, exp)

		inheritProps(&s, parent, exp)
		n.Style = s

		if len(n.Children) > 0 {
			resolveNode(&n.Children, rulesBySel, vars, &n.Style)
		}
	}
}

func parseInt(val string) (int, bool) {
	n, err := strconv.Atoi(strings.TrimSpace(val))
	if err != nil {
		return 0, false
	}
	return n, true
}

func parseFloat(val string) (float64, bool) {
	f, err := strconv.ParseFloat(strings.TrimSpace(val), 64)
	if err != nil {
		return 0, false
	}
	return f, true
}

func parseColor(val string) (tcell.Color, bool) {
	c := tcell.GetColor(val)
	if c == tcell.ColorDefault && strings.ToLower(strings.TrimSpace(val)) != "default" {
		return tcell.ColorDefault, false
	}
	return c, true
}

func parseBool(val string) (bool, bool) {
	switch strings.ToLower(strings.TrimSpace(val)) {
	case "true", "yes", "1":
		return true, true
	case "false", "no", "0":
		return false, true
	}
	return false, false
}
