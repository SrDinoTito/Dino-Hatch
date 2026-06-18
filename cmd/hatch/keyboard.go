// Manejo de eventos de teclado: navegacion, input/textarea, cierre de modal.
package main

import (
	"strings"

	"github.com/gdamore/tcell/v2"
)

// handleKeyEvent procesa teclado: navegacion, input/textarea, cierre de modal.
func handleKeyEvent(s *AppState, e *tcell.EventKey) {
	s.SelActive = false

	// C2: Tab / Shift+Tab navigation (no afecta modal, return inmediato)
	if e.Key() == tcell.KeyTab {
		s.buildFocusOrder()
		if len(s.FocusOrder) > 0 {
			if e.Modifiers()&tcell.ModShift != 0 {
				s.FocusIndex--
				if s.FocusIndex < 0 {
					s.FocusIndex = len(s.FocusOrder) - 1
				}
			} else {
				s.FocusIndex++
				if s.FocusIndex >= len(s.FocusOrder) {
					s.FocusIndex = 0
				}
			}
			s.FocusedElement = s.FocusOrder[s.FocusIndex]
			if _, ok := s.InputStates[s.FocusedElement]; !ok {
				s.InputStates[s.FocusedElement] = &inputState{}
			}
		}
		return
	}

	if s.ModalOpen {
		if e.Key() == tcell.KeyEsc {
			s.ModalOpen = false
		}
		return
	}
	if s.FocusedElement != nil && s.FocusedElement.Tag == "textarea" {
		isEnter := e.Key() == tcell.KeyEnter || e.Key() == tcell.KeyLF || e.Rune() == '\n' || e.Rune() == '\r'
		if isEnter {
			handleInputKey(s.FocusedElement, e)
			return
		}
	}
	if s.FocusedElement != nil && (s.FocusedElement.Tag == "input" || s.FocusedElement.Tag == "textarea") {
		if handleInputKey(s.FocusedElement, e) {
			return
		}
	}
	switch e.Key() {
	case tcell.KeyUp:
		s.ScrollY = max(0, s.ScrollY-1)
	case tcell.KeyDown:
		s.ScrollY = min(s.MaxScroll, s.ScrollY+1)
	case tcell.KeyPgUp:
		s.ScrollY = max(0, s.ScrollY-s.H)
	case tcell.KeyPgDn:
		s.ScrollY = min(s.MaxScroll, s.ScrollY+s.H)
	case tcell.KeyHome:
		s.ScrollY = 0
	case tcell.KeyEnd:
		s.ScrollY = s.MaxScroll
	default:
		err := s.Eng.Handle(e)
		if err != nil && strings.Contains(err.Error(), "no action") {
			// tecla no mapeada, ignorar
		}
	}
}
