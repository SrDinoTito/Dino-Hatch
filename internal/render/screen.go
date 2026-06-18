// Package render implementa el pipeline de renderizado a terminal.
// Screen abstrae la terminal real para permitir tests sin terminal fisica.
package render

import "github.com/gdamore/tcell/v2"

// Screen abstrae la terminal para poder testear sin terminal real.
type Screen interface {
	Init() error
	Flush() error
	Close() error
	Size() (int, int)
	SetContent(x, y int, ch rune, comb []rune, style tcell.Style)
}

// tcellScreen implementa Screen usando tcell real.
type tcellScreen struct {
	screen tcell.Screen
}

// NewTcellScreen crea un wrapper sobre tcell.Screen.
func NewTcellScreen() *tcellScreen {
	return &tcellScreen{}
}

// Init crea e inicializa la pantalla tcell.
func (ts *tcellScreen) Init() error {
	s, err := tcell.NewScreen()
	if err != nil {
		return err
	}
	if err := s.Init(); err != nil {
		return err
	}
	ts.screen = s
	return nil
}

// Flush envia el buffer interno a la terminal real.
func (ts *tcellScreen) Flush() error {
	ts.screen.Show()
	return nil
}

// Close finaliza la pantalla tcell (restaura la terminal).
func (ts *tcellScreen) Close() error {
	ts.screen.Fini()
	return nil
}

// Size devuelve el ancho y alto de la terminal.
func (ts *tcellScreen) Size() (int, int) {
	return ts.screen.Size()
}

// SetContent establece el contenido de una celda en (x, y).
func (ts *tcellScreen) SetContent(x, y int, ch rune, comb []rune, style tcell.Style) {
	ts.screen.SetContent(x, y, ch, comb, style)
}

// mockScreen implementa Screen para tests (guarda contenido en memoria).
type mockScreen struct {
	cells   map[[2]int]mockCell
	width   int
	height  int
	flushed bool
}

type mockCell struct {
	ch    rune
	comb  []rune
	style tcell.Style
}

// NewMockScreen crea un mock con las dimensiones dadas.
func NewMockScreen(w, h int) *mockScreen {
	return &mockScreen{
		cells:  make(map[[2]int]mockCell),
		width:  w,
		height: h,
	}
}

// Init no requiere inicializacion real en el mock.
func (ms *mockScreen) Init() error { return nil }

// Flush marca el estado como flushed.
func (ms *mockScreen) Flush() error {
	ms.flushed = true
	return nil
}

// Close no requiere accion en el mock.
func (ms *mockScreen) Close() error { return nil }

// Size devuelve las dimensiones configuradas al crear el mock.
func (ms *mockScreen) Size() (int, int) {
	return ms.width, ms.height
}

// SetContent guarda la celda en el mapa interno.
func (ms *mockScreen) SetContent(x, y int, ch rune, comb []rune, style tcell.Style) {
	ms.cells[[2]int{x, y}] = mockCell{ch: ch, comb: comb, style: style}
}

// GetCell devuelve el contenido en (x, y) — metodo auxiliar para tests.
func (ms *mockScreen) GetCell(x, y int) (rune, tcell.Style, bool) {
	cell, ok := ms.cells[[2]int{x, y}]
	if !ok {
		return 0, tcell.StyleDefault, false
	}
	return cell.ch, cell.style, true
}
