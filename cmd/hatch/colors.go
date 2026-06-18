// Modo de colores aleatorios para boxes: toggle ON/OFF colorea cada <box>
// con un color RGB unico generado aleatoriamente.
package main

import (
	"math/rand"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)

// toggleRandomColors activa o desactiva el modo de colores aleatorios.
// Al activar, recolecta todos los <box> y les asigna colores RGB unicos.
// Al desactivar, limpia el mapa de asignaciones.
func toggleRandomColors(s *AppState) {
	s.RandomColorsMode = !s.RandomColorsMode
	if s.RandomColorsMode {
		var boxes []*ast.ElementNode
		for i := range s.Doc.Pages {
			collectBoxes(&s.Doc.Pages[i].Children, &boxes)
		}
		colors := generateUniqueColors(len(boxes))
		s.BoxColors = make(map[*ast.ElementNode]tcell.Color, len(boxes))
		for j, box := range boxes {
			s.BoxColors[box] = colors[j]
		}
	} else {
		s.BoxColors = nil
	}
}

// collectBoxes recorre el arbol recursivamente y agrega punteros a elementos
// con Tag == "box". Los punteros son estables porque el AST no se reasigna.
func collectBoxes(children *[]ast.ElementNode, boxes *[]*ast.ElementNode) {
	for i := range *children {
		node := &(*children)[i]
		if node.Tag == "box" {
			*boxes = append(*boxes, node)
		}
		collectBoxes(&node.Children, boxes)
	}
}

// generateUniqueColors genera n colores RGB sin repetir.
// Usa el rango 30-230 para evitar tonos muy claros u oscuros.
func generateUniqueColors(n int) []tcell.Color {
	used := make(map[int64]bool)
	colors := make([]tcell.Color, 0, n)
	for len(colors) < n {
		r := rand.Intn(201) + 30
		g := rand.Intn(201) + 30
		b := rand.Intn(201) + 30
		key := int64(r)<<16 | int64(g)<<8 | int64(b)
		if used[key] {
			continue
		}
		used[key] = true
		colors = append(colors, tcell.NewRGBColor(int32(r), int32(g), int32(b)))
	}
	return colors
}
