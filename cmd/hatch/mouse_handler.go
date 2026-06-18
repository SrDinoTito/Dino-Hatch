// Manejo de eventos de mouse: click, drag, wheel, press.
package main
import (
	"strings"
	"github.com/gdamore/tcell/v2"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/layout"
)
func handleMouseEvent(s *AppState, e *tcell.EventMouse) {
	mx, my := e.Position()
	buttons := e.Buttons()
	py := my + s.ScrollY
	if s.ModalOpen && s.ModalDoc != nil {
		if buttons == tcell.WheelUp || buttons == tcell.WheelDown { return }
		modalEl := hitTest(s.ModalDoc, mx, my)
		if buttons != tcell.ButtonNone && s.PrevButtons == tcell.ButtonNone {
			if modalEl != nil {
				if modalEl.Tag == "button" {
					if modalEl.Attrs["id"] == "btn-cancel-modal" || modalEl.Attrs["id"] == "btn-confirm-modal" { s.ModalOpen = false }
				} else if modalEl.Tag == "box" && modalEl.Attrs["id"] == "modal-bg" { s.ModalOpen = false }
			}
		}
		s.PrevButtons = buttons; return
	}
	if buttons != tcell.ButtonNone && buttons != tcell.WheelUp && buttons != tcell.WheelDown && s.PrevButtons != tcell.ButtonNone {
		s.SelEndX = mx; s.SelEndY = py; s.SelActive = true
	}
	if buttons == tcell.WheelUp || buttons == tcell.WheelDown {
		if s.FocusedElement != nil && s.FocusedElement.Tag == "textarea" {
			bb := s.FocusedElement.BoundBox
			if mx >= bb.X && mx < bb.X+bb.W && my >= bb.Y && my < bb.Y+bb.H {
				st := s.InputStates[s.FocusedElement]
				if st == nil { st = &inputState{BaseLines: 1}; s.InputStates[s.FocusedElement] = st }
				lines := strings.Split(s.FocusedElement.Text, "\n")
				maxScroll := len(lines) - (st.BaseLines + 4)
				if maxScroll < 0 { maxScroll = 0 }
				if buttons == tcell.WheelUp { if st.ScrollY > 0 { st.ScrollY-- } } else { if st.ScrollY < maxScroll { st.ScrollY++ } }
				s.PrevButtons = buttons; return
			}
		}
		// D3: scroll container handling — scrolling interno en vez de pagina
		scrolled := false
		for i := range s.Doc.Pages {
			if el := findScrollContainer(&s.Doc.Pages[i], mx, my, s.ScrollY); el != nil {
				contentH := layout.ContentHeight(el)
				maxScroll := max(0, contentH-el.BoundBox.H)
				if buttons == tcell.WheelUp && el.ScrollY > 0 {
					el.ScrollY -= 3
					s.Dirty = true
				} else if buttons == tcell.WheelDown && el.ScrollY < maxScroll {
					el.ScrollY += 3
					s.Dirty = true
				}
				scrolled = true
				break
			}
		}
		if !scrolled {
			// scroll normal de pagina
			if buttons == tcell.WheelUp { s.ScrollY = max(0, s.ScrollY-3) } else { s.ScrollY = min(s.MaxScroll, s.ScrollY+3) }
		}
		s.PrevButtons = buttons; return
	}
	el := hitTest(s.Doc, mx, py)
	handleMousePress(s, mx, my, py, buttons, el)
	handleMouseRelease(s, buttons, el)
	s.HoveredElement = el; s.PrevButtons = buttons
}
func handleMousePress(s *AppState, mx, my, py int, buttons tcell.ButtonMask, el *ast.ElementNode) {
	if buttons != tcell.ButtonNone && s.PrevButtons == tcell.ButtonNone {
		s.SelStartX, s.SelStartY = mx, py
		s.SelEndX, s.SelEndY = mx, py
		s.SelActive = false
		if el != nil && el.Tag == "button" && el.Attrs["id"] == "open-modal" { s.ModalOpen = true; s.FocusedElement = nil; s.PrevButtons = buttons; return }
		if el != nil && el.Tag == "button" {
			if id, ok := el.Attrs["id"]; ok && strings.HasPrefix(id, "nav-") {
				if _, exists := s.PageFiles[id[4:]]; exists { s.LoadPage(id[4:]); s.CurrentPage = id[4:]; s.FocusedElement = nil; s.ScrollY = 0; s.Dirty = true; s.PrevButtons = buttons; return }
			}
		}
		if el != nil && el.Events["click"] != "" {
			executeEvent(s, el, "click")
		}
		if el != nil && (el.Tag == "input" || el.Tag == "textarea" || el.Tag == "button") {
			s.FocusedElement = el
			if el.Tag == "input" || el.Tag == "textarea" {
				st := s.InputStates[el]
				if st == nil { st = &inputState{}; s.InputStates[el] = st }
				relX := mx - el.BoundBox.X
				if el.Style.Border { relX-- }
				if el.Tag == "textarea" {
					relY := my - el.BoundBox.Y
					if el.Style.Border { relY-- }
					if relY < 0 { relY = 0 }
					lines := strings.Split(el.Text, "\n")
					if relY >= len(lines) { relY = len(lines) - 1 }
					pos := 0
					for i := 0; i < relY; i++ { pos += len([]rune(lines[i])) + 1 }
					if relX < 0 { relX = 0 }
					if textLen := len([]rune(lines[relY])); relX > textLen { relX = textLen }
					st.Cursor = pos + relX
				} else {
					if relX < 0 { relX = 0 }
					if textLen := len([]rune(el.Text)); relX > textLen { relX = textLen }
					st.Cursor = relX
				}
			}
		} else { s.FocusedElement = nil }
	}
}
