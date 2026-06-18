// Tests edge case para CellBuffer.
// Cubre Diff con tamanos diferentes, Get fuera de rango,
// acceso via Cells(), y Fill con estilo especifico.
package render

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// TestDiff_DifferentSizes verifica que Diff entre buffers de diferente
// tamano compare solo el area comun y no haga panic.
func TestDiff_DifferentSizes(t *testing.T) {
	cb1 := NewCellBuffer(10, 5)
	cb2 := NewCellBuffer(5, 10)
	cb1.Set(4, 4, 'A', tcell.StyleDefault)
	cb2.Set(4, 4, 'B', tcell.StyleDefault)
	rect, updates := cb1.Diff(cb2)
	// Area comun: 5x5 → (4,4) esta dentro
	if len(updates) != 1 {
		t.Errorf("esperado 1 update en area comun, got %d", len(updates))
	}
	_ = rect
	_ = updates
}

// TestDiff_DifferentSizes_NoOverlap verifica que cambios fuera del area
// comun No se reporten como updates.
func TestDiff_DifferentSizes_NoOverlap(t *testing.T) {
	// cb1=3x3, cb2=5x3, diferencias solo en x>=3 (fuera de area comun)
	cb1 := NewCellBuffer(3, 3)
	cb2 := NewCellBuffer(5, 3)
	cb1.Fill('.', tcell.StyleDefault)
	cb2.Fill('.', tcell.StyleDefault)
	// Cambiar celda en x=4 (fuera del area comun 0..2)
	cb2.Set(4, 1, 'X', tcell.StyleDefault)
	rect, updates := cb1.Diff(cb2)
	if len(updates) != 0 {
		t.Errorf("esperado 0 updates (cambio fuera de area comun), got %d", len(updates))
	}
	if rect.W != 0 || rect.H != 0 {
		t.Errorf("esperado dirty rect vacio, got %+v", rect)
	}
}

// TestGetOutOfBounds verifica que Get devuelva error con coordenadas
// fuera de rango: x negativo, x igual a width, y igual a height.
func TestGetOutOfBounds(t *testing.T) {
	cb := NewCellBuffer(5, 3)
	_, err := cb.Get(-1, 0)
	if err == nil {
		t.Error("esperado error por x negativo")
	}
	_, err = cb.Get(5, 0)
	if err == nil {
		t.Error("esperado error por x=width")
	}
	_, err = cb.Get(0, 3)
	if err == nil {
		t.Error("esperado error por y=height")
	}
}

// TestCells verifica que Cells() devuelva el slice interno con
// longitud = w*h, y que modificarlo via indice directo se refleje en Get.
func TestCells(t *testing.T) {
	cb := NewCellBuffer(4, 3)
	cells := cb.Cells()
	if len(cells) != 12 {
		t.Errorf("esperado 12 celdas (4x3), got %d", len(cells))
	}
	// Modificar via slice directo debe reflejarse en Get
	cells[0] = Cell{Rune: 'Z', Style: tcell.StyleDefault}
	cell, _ := cb.Get(0, 0)
	if cell.Rune != 'Z' {
		t.Errorf("esperado 'Z' via Cells(), got %c", cell.Rune)
	}
}

// TestFill_StylePersist verifica que Fill con estilo especifico
// aplique el mismo rune y estilo a todas las celdas.
func TestFill_StylePersist(t *testing.T) {
	style := tcell.StyleDefault.Foreground(tcell.ColorBlue).Background(tcell.ColorWhite)
	cb := NewCellBuffer(2, 2)
	cb.Fill('X', style)
	for y := 0; y < 2; y++ {
		for x := 0; x < 2; x++ {
			cell, _ := cb.Get(x, y)
			if cell.Rune != 'X' {
				t.Errorf("esperado 'X' en (%d,%d), got %c", x, y, cell.Rune)
			}
			if cell.Style != style {
				t.Errorf("estilo no coincide en (%d,%d)", x, y)
			}
		}
	}
}
