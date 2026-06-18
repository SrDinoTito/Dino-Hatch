// Package render implementa el pipeline de renderizado a terminal.
// CellBuffer es el buffer intermedio para diff entre frames.
package render

import (
	"fmt"

	"github.com/gdamore/tcell/v2"
)

// Cell representa una celda individual en el buffer de pantalla.
type Cell struct {
	Rune  rune
	Style tcell.Style
}

// DirtyRect es el rectángulo que engloba todas las celdas modificadas.
// Si W y H son 0, no hay cambios.
type DirtyRect struct {
	X, Y, W, H int
}

// CellUpdate representa una celda que cambió entre dos frames.
type CellUpdate struct {
	X, Y int
	Cell Cell
}

// CellBuffer es una matriz de celdas (W x H) en orden row-major.
// Almacena el estado de la pantalla para comparar entre frames y
// minimizar actualizaciones al driver de terminal.
type CellBuffer struct {
	cells []Cell
	w, h  int
}

// NewCellBuffer crea un nuevo buffer con las dimensiones dadas.
func NewCellBuffer(w, h int) *CellBuffer {
	return &CellBuffer{
		cells: make([]Cell, w*h),
		w:     w,
		h:     h,
	}
}

// Set establece el contenido de una celda en (x, y).
func (cb *CellBuffer) Set(x, y int, r rune, style tcell.Style) error {
	if x < 0 || x >= cb.w || y < 0 || y >= cb.h {
		return fmt.Errorf(
			"cellbuffer: posicion (%d,%d) fuera de rango para buffer %dx%d",
			x, y, cb.w, cb.h,
		)
	}
	cb.cells[y*cb.w+x] = Cell{Rune: r, Style: style}
	return nil
}

// Get obtiene el contenido de una celda en (x, y).
func (cb *CellBuffer) Get(x, y int) (Cell, error) {
	if x < 0 || x >= cb.w || y < 0 || y >= cb.h {
		return Cell{}, fmt.Errorf(
			"cellbuffer: posicion (%d,%d) fuera de rango para buffer %dx%d",
			x, y, cb.w, cb.h,
		)
	}
	return cb.cells[y*cb.w+x], nil
}

// Fill llena todo el buffer con un rune y estilo.
func (cb *CellBuffer) Fill(r rune, style tcell.Style) {
	c := Cell{Rune: r, Style: style}
	for i := range cb.cells {
		cb.cells[i] = c
	}
}

// Diff compara este buffer con otro y devuelve las celdas cambiadas.
// Ambos buffers deben tener el mismo tamano; si no, se compara el area comun.
func (cb *CellBuffer) Diff(other *CellBuffer) (DirtyRect, []CellUpdate) {
	w := cb.w
	h := cb.h
	if other.w < w {
		w = other.w
	}
	if other.h < h {
		h = other.h
	}

	var updates []CellUpdate
	minX, minY := w, h
	maxX, maxY := -1, -1

	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			idx := y*cb.w + x
			if cb.cells[idx] != other.cells[y*other.w+x] {
				updates = append(updates, CellUpdate{X: x, Y: y, Cell: other.cells[y*other.w+x]})
				if x < minX {
					minX = x
				}
				if y < minY {
					minY = y
				}
				if x > maxX {
					maxX = x
				}
				if y > maxY {
					maxY = y
				}
			}
		}
	}

	if len(updates) == 0 {
		return DirtyRect{}, nil
	}

	return DirtyRect{X: minX, Y: minY, W: maxX - minX + 1, H: maxY - minY + 1}, updates
}

// Resize cambia el tamano del buffer y reinicializa el contenido.
func (cb *CellBuffer) Resize(w, h int) {
	cb.cells = make([]Cell, w*h)
	cb.w = w
	cb.h = h
}

// Width devuelve el ancho del buffer.
func (cb *CellBuffer) Width() int {
	return cb.w
}

// Height devuelve el alto del buffer.
func (cb *CellBuffer) Height() int {
	return cb.h
}

// Cells devuelve el slice interno de celdas para acceso directo.
func (cb *CellBuffer) Cells() []Cell {
	return cb.cells
}
