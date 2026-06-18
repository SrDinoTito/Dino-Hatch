package main
import (
	"strings"

	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
)
func cursorLineCol(text string, cursor int) (int, int) {
	lines := strings.Split(text, "\n")
	if cursor > len([]rune(text)) {
		cursor = len([]rune(text))
	}
	pos := 0
	for li, l := range lines {
		if pos+len(l) >= cursor {
			col := cursor - pos
			if col > len(l) {
				col = len(l)
			}
			return li, col
		}
		pos += len(l) + 1
		if li == len(lines)-1 {
			return li, len(l)
		}
	}
	return 0, 0
}
func moveCursorUp(el *ast.ElementNode, st *inputState) bool {
	if el.Tag != "textarea" {
		return false
	}
	lines := strings.Split(el.Text, "\n")
	cursorLine, cursorCol := cursorLineCol(el.Text, st.Cursor)
	if cursorLine > 0 {
		prevLen := len(lines[cursorLine-1])
		if cursorCol > prevLen {
			cursorCol = prevLen
		}
		newPos := 0
		for i := 0; i < cursorLine-1; i++ {
			newPos += len(lines[i]) + 1
		}
		st.Cursor = newPos + cursorCol
	}
	return true
}
func moveCursorDown(el *ast.ElementNode, st *inputState) bool {
	if el.Tag != "textarea" {
		return false
	}
	lines := strings.Split(el.Text, "\n")
	cursorLine, cursorCol := cursorLineCol(el.Text, st.Cursor)
	if cursorLine < len(lines)-1 {
		nextLen := len(lines[cursorLine+1])
		if cursorCol > nextLen {
			cursorCol = nextLen
		}
		newPos := 0
		for i := 0; i <= cursorLine; i++ {
			newPos += len(lines[i]) + 1
		}
		st.Cursor = newPos + cursorCol
	}
	return true
}
func handleInputKey(el *ast.ElementNode, ev *tcell.EventKey) bool {
	st := state.InputStates[el]
	if st == nil {
		st = &inputState{Cursor: len([]rune(el.Text))}
		state.InputStates[el] = st
	}
	switch ev.Key() {
	case tcell.KeyRune:
		r := ev.Rune()
		text := []rune(el.Text)
		pos := st.Cursor
		if pos > len(text) {
			pos = len(text)
		}
		if (r == '\r' || r == '\n') && el.Tag == "textarea" {
			newText := string(text[:pos]) + "\n" + string(text[pos:])
			el.Text = newText
			state.LayoutDirty = true
			st.Cursor = pos + 1
			executeDataBinding(state, el)
		} else {
			newText := string(text[:pos]) + string(r) + string(text[pos:])
			el.Text = newText
			state.LayoutDirty = true
			st.Cursor = pos + 1
			executeDataBinding(state, el)
		}
		return true
	case tcell.KeyBackspace, tcell.KeyBackspace2:
		if st.Cursor > 0 {
			text := []rune(el.Text)
			newText := string(text[:st.Cursor-1]) + string(text[st.Cursor:])
			el.Text = newText
			state.LayoutDirty = true
			st.Cursor--
			executeDataBinding(state, el)
		}
		return true
	case tcell.KeyDelete:
		text := []rune(el.Text)
		if st.Cursor < len(text) {
			newText := string(text[:st.Cursor]) + string(text[st.Cursor+1:])
			el.Text = newText
			state.LayoutDirty = true
			executeDataBinding(state, el)
		}
		return true
	case tcell.KeyLeft:
		if st.Cursor > 0 {
			st.Cursor--
		}
		return true
	case tcell.KeyRight:
		text := []rune(el.Text)
		if st.Cursor < len(text) {
			st.Cursor++
		}
		return true
	case tcell.KeyUp:
		return moveCursorUp(el, st)
	case tcell.KeyDown:
		return moveCursorDown(el, st)
	case tcell.KeyEnter, tcell.KeyLF:
		if el.Tag == "textarea" {
			text := []rune(el.Text)
			pos := st.Cursor
			if pos > len(text) {
				pos = len(text)
			}
			newText := string(text[:pos]) + "\n" + string(text[pos:])
			el.Text = newText
			state.LayoutDirty = true
			st.Cursor = pos + 1
			executeDataBinding(state, el)
			return true
		}
		return false
	case tcell.KeyEsc:
		state.FocusedElement = nil
		return true
	default:
		return false
	}
}
