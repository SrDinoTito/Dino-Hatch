package main

import (
	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/render"
)

func renderBoxContent(cb *render.CellBuffer, node *ast.ElementNode, style tcell.Style, x, y, w, h int) {
	if node.Style.Border {
		drawBorder(cb, x, y, w, h, style)
		fillRect(cb, x+1, y+1, w-2, h-2, ' ', style)
	} else {
		fillRect(cb, x, y, w, h, ' ', style)
	}
}

func renderButton(cb *render.CellBuffer, node *ast.ElementNode, style tcell.Style, bg tcell.Color, x, y, w, h, scrollY int) {
	if node == state.HoveredElement {
		hoverBg := bg
		if hoverBg == tcell.ColorReset {
			hoverBg = tcell.ColorGray
		}
		bg = hoverBg
		style = style.Background(hoverBg)
	}
	if node.Style.Border {
		drawBorder(cb, x, y, w, h, style)
		fillRect(cb, x+1, y+1, w-2, h-2, ' ', style)
	} else {
		fillRect(cb, x, y, w, h, ' ', style)
	}
	innerX, innerY, innerW, innerH := x, y, w, h
	if node.Style.Border {
		innerX = x + 1
		innerY = y + 1
		innerW = w - 2
		innerH = h - 2
	}
	for i := range node.Children {
		child := &node.Children[i]
		if child.Tag == "text" {
			textStyle := style.Foreground(child.Style.Color)
			xOffset := 0
			switch child.Style.TextAlign {
			case "center":
				xOffset = max(0, (innerW-len(child.Text))/2)
			case "right":
				xOffset = max(0, innerW-len(child.Text))
			}
			yOffset := 0
			switch child.Style.VAlign {
			case "middle":
				yOffset = max(0, (innerH-1)/2)
			case "bottom":
				yOffset = max(0, innerH-1)
			}
			for j, r := range child.Text {
				tx := innerX + xOffset + j
				ty := innerY + yOffset
				if tx >= innerX+innerW || tx >= cb.Width() || ty >= innerY+innerH || ty >= cb.Height() || tx < 0 || ty < 0 {
					continue
				}
				_ = cb.Set(tx, ty, r, textStyle)
			}
		} else {
			renderNode(cb, child, bg, scrollY)
		}
	}
}

func renderText(cb *render.CellBuffer, node *ast.ElementNode, style tcell.Style, x, y, w, h int) {
	xOffset := 0
	switch node.Style.TextAlign {
	case "center":
		xOffset = max(0, (w-len(node.Text))/2)
	case "right":
		xOffset = max(0, w-len(node.Text))
	}
	yOffset := 0
	switch node.Style.VAlign {
	case "middle":
		yOffset = max(0, (h-1)/2)
	case "bottom":
		yOffset = max(0, h-1)
	}
	selMinX, selMinY, selMaxX, selMaxY := 0, 0, -1, -1
	if state.SelActive {
		selMinX, selMinY, selMaxX, selMaxY = normalizedSelectionRect()
	}
	for i, r := range node.Text {
		tx := x + xOffset + i
		ty := y + yOffset
		if tx >= x+w || tx >= cb.Width() || ty >= y+h || ty >= cb.Height() || tx < 0 || ty < 0 {
			continue
		}
		charStyle := style
		if state.SelActive && tx >= selMinX && tx <= selMaxX && ty >= selMinY && ty <= selMaxY {
			charStyle = style.Reverse(true)
		}
		_ = cb.Set(tx, ty, r, charStyle)
	}
}

func renderInput(cb *render.CellBuffer, node *ast.ElementNode, style tcell.Style, x, y, w, h int) {
	if node.Style.Border {
		drawBorder(cb, x, y, w, h, style)
		fillRect(cb, x+1, y+1, w-2, h-2, ' ', style)
	} else {
		fillRect(cb, x, y, w, h, ' ', style)
	}
	val := node.Text
	if val == "" {
		if v, ok := node.Attrs["value"]; ok {
			val = v
		}
	}
	textStyle := style.Foreground(node.Style.Color)
	tx, ty, maxW := x, y, w
	if node.Style.Border {
		tx = x + 1
		ty = y + 1
		maxW = w - 2
	}
	for i, ch := range val {
		if i >= maxW {
			break
		}
		_ = cb.Set(tx+i, ty, ch, textStyle)
	}
	if node == state.FocusedElement {
		st := state.InputStates[node]
		if st == nil {
			st = &inputState{Cursor: len([]rune(val))}
			state.InputStates[node] = st
		}
		curX := tx + st.Cursor
		if curX < tx+maxW && curX >= tx {
			_ = cb.Set(curX, ty, '▌', textStyle)
		}
	}
}

func fillRect(cb *render.CellBuffer, x, y, w, h int, ch rune, style tcell.Style) {
	for dy := 0; dy < h; dy++ {
		for dx := 0; dx < w; dx++ {
			_ = cb.Set(x+dx, y+dy, ch, style)
		}
	}
}
