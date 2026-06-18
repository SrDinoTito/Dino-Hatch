// Estado de cursor, seleccion y elementos interactivos.
package main

// inputState almacena el estado de cursor y scroll de un input o textarea.
type inputState struct {
	Cursor    int // posicion del cursor en runas dentro del texto
	ScrollY   int // desplazamiento interno para textarea
	BaseLines int // lineas iniciales del textarea, para calcular limite de expansion
}

// normalizedSelectionRect devuelve el rectangulo de seleccion normalizado
// (minX, minY, maxX, maxY) sin importar la direccion del drag.
func normalizedSelectionRect() (int, int, int, int) {
	minX := state.SelStartX
	maxX := state.SelEndX
	if state.SelEndX < state.SelStartX {
		minX = state.SelEndX
		maxX = state.SelStartX
	}
	minY := state.SelStartY
	maxY := state.SelEndY
	if state.SelEndY < state.SelStartY {
		minY = state.SelEndY
		maxY = state.SelStartY
	}
	return minX, minY, maxX, maxY
}
