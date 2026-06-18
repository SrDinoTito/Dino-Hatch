package render

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNewCellBuffer(t *testing.T) {
	cb := NewCellBuffer(10, 5)
	if cb.Width() != 10 {
		t.Errorf("esperado ancho 10, got %d", cb.Width())
	}
	if cb.Height() != 5 {
		t.Errorf("esperado alto 5, got %d", cb.Height())
	}
}

func TestSetGet(t *testing.T) {
	cb := NewCellBuffer(10, 5)
	style := tcell.StyleDefault
	err := cb.Set(3, 2, 'A', style)
	if err != nil {
		t.Fatalf("Set fallo: %v", err)
	}
	cell, err := cb.Get(3, 2)
	if err != nil {
		t.Fatalf("Get fallo: %v", err)
	}
	if cell.Rune != 'A' {
		t.Errorf("esperado 'A', got %c", cell.Rune)
	}
	if cell.Style != style {
		t.Errorf("estilo no coincide")
	}
}

func TestSetOutOfBounds(t *testing.T) {
	cb := NewCellBuffer(10, 5)
	err := cb.Set(10, 0, 'X', tcell.StyleDefault)
	if err == nil {
		t.Error("esperado error por x fuera de rango")
	}
	err = cb.Set(-1, 0, 'X', tcell.StyleDefault)
	if err == nil {
		t.Error("esperado error por x negativo")
	}
	err = cb.Set(0, 5, 'X', tcell.StyleDefault)
	if err == nil {
		t.Error("esperado error por y fuera de rango")
	}
}

func TestFill(t *testing.T) {
	cb := NewCellBuffer(3, 3)
	style := tcell.StyleDefault.Foreground(tcell.ColorRed)
	cb.Fill('\u2588', style)
	for y := 0; y < 3; y++ {
		for x := 0; x < 3; x++ {
			cell, err := cb.Get(x, y)
			if err != nil {
				t.Fatalf("Get fallo en (%d,%d): %v", x, y, err)
			}
			if cell.Rune != '\u2588' {
				t.Errorf("esperado bloque en (%d,%d), got %c", x, y, cell.Rune)
			}
			if cell.Style != style {
				t.Errorf("estilo no coincide en (%d,%d)", x, y)
			}
		}
	}
}

func TestDiff_NoChanges(t *testing.T) {
	cb1 := NewCellBuffer(4, 3)
	cb2 := NewCellBuffer(4, 3)
	cb1.Fill('.', tcell.StyleDefault)
	cb2.Fill('.', tcell.StyleDefault)
	rect, updates := cb1.Diff(cb2)
	if rect.W != 0 || rect.H != 0 {
		t.Errorf("esperado rect vacio, got %+v", rect)
	}
	if len(updates) != 0 {
		t.Errorf("esperado 0 updates, got %d", len(updates))
	}
}

func TestDiff_OneCell(t *testing.T) {
	cb1 := NewCellBuffer(4, 3)
	cb2 := NewCellBuffer(4, 3)
	cb1.Fill('.', tcell.StyleDefault)
	cb2.Fill('.', tcell.StyleDefault)
	cb2.Set(2, 1, 'X', tcell.StyleDefault)
	rect, updates := cb1.Diff(cb2)
	if len(updates) != 1 {
		t.Fatalf("esperado 1 update, got %d", len(updates))
	}
	if updates[0].X != 2 || updates[0].Y != 1 {
		t.Errorf("esperado (2,1), got (%d,%d)", updates[0].X, updates[0].Y)
	}
	if updates[0].Cell.Rune != 'X' {
		t.Errorf("esperado 'X', got %c", updates[0].Cell.Rune)
	}
	if rect.W != 1 || rect.H != 1 {
		t.Errorf("esperado rect 1x1, got %+v", rect)
	}
}
func TestDiff_DirtyRect(t *testing.T) {
	cb1 := NewCellBuffer(6, 4)
	cb2 := NewCellBuffer(6, 4)
	cb1.Fill('.', tcell.StyleDefault)
	cb2.Fill('.', tcell.StyleDefault)
	cb2.Set(2, 0, 'A', tcell.StyleDefault)
	cb2.Set(4, 2, 'B', tcell.StyleDefault)
	rect, updates := cb1.Diff(cb2)
	if len(updates) != 2 {
		t.Fatalf("esperado 2 updates, got %d", len(updates))
	}
	// DirtyRect debe ser el bounding box: (2,0) a (4,2) → {2,0,3,3}
	if rect.X != 2 || rect.Y != 0 || rect.W != 3 || rect.H != 3 {
		t.Errorf("esperado rect {2,0,3,3}, got %+v", rect)
	}
}

func TestResize(t *testing.T) {
	cb := NewCellBuffer(10, 5)
	cb.Fill('X', tcell.StyleDefault)
	cb.Resize(3, 3)
	if cb.Width() != 3 {
		t.Errorf("esperado ancho 3 tras resize, got %d", cb.Width())
	}
	if cb.Height() != 3 {
		t.Errorf("esperado alto 3 tras resize, got %d", cb.Height())
	}
	// Debe estar reinicializado, no debe conservar 'X'
	cell, _ := cb.Get(0, 0)
	if cell.Rune != 0 {
		t.Errorf("esperado rune cero tras resize, got %c", cell.Rune)
	}
}

func TestWidthHeight(t *testing.T) {
	cb := NewCellBuffer(7, 11)
	if cb.Width() != 7 {
		t.Errorf("esperado 7, got %d", cb.Width())
	}
	if cb.Height() != 11 {
		t.Errorf("esperado 11, got %d", cb.Height())
	}
}
