package render

import (
	"testing"

	"github.com/gdamore/tcell/v2"
)

// TestTcellScreen_ImplementsInterface verifica que *tcellScreen cumple Screen.
func TestTcellScreen_ImplementsInterface(t *testing.T) {
	var _ Screen = (*tcellScreen)(nil)
}

// TestMockScreen_ImplementsInterface verifica que *mockScreen cumple Screen.
func TestMockScreen_ImplementsInterface(t *testing.T) {
	var _ Screen = (*mockScreen)(nil)
}

// TestMockScreen_SetGetContent verifica set/get de celdas en mock.
func TestMockScreen_SetGetContent(t *testing.T) {
	ms := NewMockScreen(80, 24)
	ms.Init()
	ms.SetContent(10, 5, 'A', nil, tcell.StyleDefault)
	ch, style, ok := ms.GetCell(10, 5)
	if !ok {
		t.Error("GetCell devolvio false para celda existente")
	}
	if ch != 'A' {
		t.Errorf("esperado 'A', got %c", ch)
	}
	if style != tcell.StyleDefault {
		t.Error("estilo no coincide")
	}
	// Celda no asignada debe devolver false
	_, _, ok = ms.GetCell(0, 0)
	if ok {
		t.Error("GetCell devolvio true para celda no existente")
	}
}

// TestMockScreen_Flush verifica que Flush marca el estado.
func TestMockScreen_Flush(t *testing.T) {
	ms := NewMockScreen(80, 24)
	ms.Init()
	if ms.flushed {
		t.Error("flushed debe ser false antes de Flush")
	}
	ms.Flush()
	if !ms.flushed {
		t.Error("flushed debe ser true tras Flush")
	}
}

// TestMockScreen_Size verifica que las dimensiones sean correctas.
func TestMockScreen_Size(t *testing.T) {
	ms := NewMockScreen(120, 30)
	w, h := ms.Size()
	if w != 120 {
		t.Errorf("esperado ancho 120, got %d", w)
	}
	if h != 30 {
		t.Errorf("esperado alto 30, got %d", h)
	}
}

// TestMockScreen_Close verifica que Close no devuelva error.
func TestMockScreen_Close(t *testing.T) {
	ms := NewMockScreen(80, 24)
	ms.Init()
	err := ms.Close()
	if err != nil {
		t.Errorf("Close devolvio error inesperado: %v", err)
	}
}

// TestNewTcellScreen_ReturnsNonNil verifica que NewTcellScreen devuelva un puntero no nil.
func TestNewTcellScreen_ReturnsNonNil(t *testing.T) {
	ts := NewTcellScreen()
	if ts == nil {
		t.Error("NewTcellScreen devolvio nil")
	}
}
