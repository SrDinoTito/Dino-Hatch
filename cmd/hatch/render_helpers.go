package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/render"
)

func drawBorder(cb *render.CellBuffer, x, y, w, h int, style tcell.Style) {
	if w <= 1 || h <= 1 {
		return
	}
	for dx := 1; dx < w-1; dx++ {
		_ = cb.Set(x+dx, y, '─', style)
		_ = cb.Set(x+dx, y+h-1, '─', style)
	}
	for dy := 1; dy < h-1; dy++ {
		_ = cb.Set(x, y+dy, '│', style)
		_ = cb.Set(x+w-1, y+dy, '│', style)
	}
	_ = cb.Set(x, y, '┌', style)
	_ = cb.Set(x+w-1, y, '┐', style)
	_ = cb.Set(x, y+h-1, '└', style)
	_ = cb.Set(x+w-1, y+h-1, '┘', style)
}

func drawOverlay(cb *render.CellBuffer) {
	for y := 0; y < cb.Height(); y++ {
		for x := 0; x < cb.Width(); x++ {
			cell, err := cb.Get(x, y)
			if err == nil {
				_, bg, _ := cell.Style.Decompose()
				cb.Set(x, y, cell.Rune, cell.Style.Dim(true).Background(bg))
			}
		}
	}
}

func renderTextarea(cb *render.CellBuffer, node *ast.ElementNode, style tcell.Style, x, y, w, h int) {
	if node.Style.Border {
		drawBorder(cb, x, y, w, h, style)
		fillRect(cb, x+1, y+1, w-2, h-2, ' ', style)
	} else {
		fillRect(cb, x, y, w, h, ' ', style)
	}
	tval := node.Text
	if tval == "" {
		if v, ok := node.Attrs["value"]; ok {
			tval = v
		}
	}
	lines := strings.Split(tval, "\n")
	tts := style.Foreground(node.Style.Color)
	ttx, tmaxW, tmaxH := x, w, h
	if node.Style.Border {
		ttx = x + 1
		tmaxW = w - 2
		tmaxH = h - 2
	}
	st := state.InputStates[node]
	if st == nil {
		st = &inputState{BaseLines: len(lines)}
		if st.BaseLines < 1 {
			st.BaseLines = 1
		}
		if node.Style.MaxHeight == 0 {
			node.Style.MaxHeight = st.BaseLines + 4 + 2
		}
		state.InputStates[node] = st
	}
	maxVisible := st.BaseLines + 4
	if maxVisible > tmaxH {
		maxVisible = tmaxH
	}
	for lineIdx := st.ScrollY; lineIdx < len(lines) && lineIdx < st.ScrollY+maxVisible; lineIdx++ {
		sy := y + (lineIdx - st.ScrollY)
		if node.Style.Border {
			sy = y + 1 + (lineIdx - st.ScrollY)
		}
		for i, ch := range lines[lineIdx] {
			if i >= tmaxW {
				break
			}
			_ = cb.Set(ttx+i, sy, ch, tts)
		}
	}
	if node == state.FocusedElement {
		if st.Cursor > len([]rune(tval)) {
			st.Cursor = len([]rune(tval))
		}
		cursorLine, cursorCol := 0, st.Cursor
		pos := 0
		for li, l := range lines {
			if pos+len(l) >= st.Cursor {
				cursorLine = li
				cursorCol = st.Cursor - pos
				break
			}
			pos += len(l) + 1
			if li+1 >= len(lines) {
				cursorLine = li + 1
				cursorCol = st.Cursor - pos
			}
		}
		if cursorLine < st.ScrollY {
			st.ScrollY = cursorLine
		} else if cursorLine >= st.ScrollY+maxVisible {
			st.ScrollY = cursorLine - maxVisible + 1
		}
		visualLine := cursorLine - st.ScrollY
		if visualLine < maxVisible && visualLine >= 0 {
			if cursorCol > tmaxW {
				cursorCol = tmaxW
			}
			tsy := y + visualLine
			if node.Style.Border {
				tsy = y + 1 + visualLine
			}
			curX := ttx + cursorCol
			if curX < ttx+tmaxW && curX >= ttx {
				_ = cb.Set(curX, tsy, '▌', tts)
			}
		}
	}
	// Wave 5: Scrollbar interno del textarea
	renderTextareaScrollbar(cb, node, st, x, y, w, h)
}

func renderStickyNavbar(cb *render.CellBuffer, doc *ast.Document, w int) {
	if navbarEl := findElementByID(doc, "navbar"); navbarEl != nil {
		renderNode(cb, navbarEl, tcell.ColorReset, 0)
	}
}

func renderModal(cb *render.CellBuffer, doc *ast.Document, screenW, screenH int) {
	drawOverlay(cb)
	if len(doc.Pages) > 0 && len(doc.Pages[0].Children) > 0 {
		bg := &doc.Pages[0].Children[0]
		for i := range bg.Children {
			renderNode(cb, &bg.Children[i], tcell.ColorReset, 0)
		}
	}
}
