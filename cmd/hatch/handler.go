// Manejo de eventos declarativos (onclick, onchange) y data binding (bind).
package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"strings"

	"github.com/srdino/dino-hatch/internal/actions"
	"github.com/srdino/dino-hatch/internal/ast"
	"github.com/srdino/dino-hatch/internal/parser"
)

var actionRegistry = actions.New()

// initActions crea los callbacks y registra los handlers de acciones.
func initActions(s *AppState) {
	cb := &actions.Callbacks{
		Navigate: func(page string) {
			if _, exists := s.PageFiles[page]; exists {
				s.LoadPage(page)
				s.CurrentPage = page
				s.FocusedElement = nil
				s.ScrollY = 0
				s.Dirty = true
			}
		},
		ModalOpen:   func() { s.ModalOpen = true; s.FocusedElement = nil; s.Dirty = true },
		ModalClose:  func() { s.ModalOpen = false; s.Dirty = true },
		ModalToggle: func() { s.ModalOpen = !s.ModalOpen; s.Dirty = true },
		RandomColors: func() { toggleRandomColors(s); s.Dirty = true },
		Quit:        func() { s.Running = false },
		SetDirty:    func() { s.Dirty = true },
		ExecOutput: func(cmd string) (string, error) {
			parts := strings.Fields(cmd)
			if len(parts) == 0 {
				return "", nil
			}
			out, err := exec.Command(parts[0], parts[1:]...).CombinedOutput()
			return string(out), err
		},
		HTTPGet: func(url string) (string, error) {
			resp, err := http.Get(url)
			if err != nil {
				return "", err
			}
			defer resp.Body.Close()
			body, _ := io.ReadAll(resp.Body)
			return string(body), nil
		},
		SetElementText: func(id, text string) {
			el := s.FindElementByID(id)
			if el != nil {
				el.Text += text
				s.LayoutDirty = true
				s.Dirty = true
			}
		},
		Log: func(format string, args ...interface{}) {
			log.Printf(format, args...)
		},
		FindElement: func(id string) interface{} {
			return s.FindElementByID(id)
		},
		PublishEvent: func(topic, data string) {
			if s.EventBus != nil {
				s.EventBus.Publish(topic, data)
			}
		},
		SwitchTheme: func(name string) error {
			if s.ThemeManager == nil {
				return fmt.Errorf("no hay ThemeManager")
			}
			vars, err := s.ThemeManager.Switch(name)
			if err != nil {
				return err
			}
			if s.Doc != nil {
				s.Doc.ThemeVars = vars
				parser.ComputeStyles(s.Doc, s.StyleRules, vars)
				s.LayoutDirty = true
				s.Dirty = true
			}
			return nil
		},
	}
	s.ActionsCB = cb

	// Registrar handlers built-in
	actionRegistry.Register("page", func(args string, cb *actions.Callbacks) error {
		cb.Navigate(args)
		return nil
	})
	actionRegistry.Register("modal", func(args string, cb *actions.Callbacks) error {
		switch args {
		case "open":
			cb.ModalOpen()
		case "close":
			cb.ModalClose()
		case "toggle":
			cb.ModalToggle()
		}
		return nil
	})
	actionRegistry.Register("action", func(args string, cb *actions.Callbacks) error {
		switch args {
		case "random_colors":
			cb.RandomColors()
		case "quit":
			cb.Quit()
		default:
			return fmt.Errorf("accion desconocida: %s", args)
		}
		return nil
	})
	actionRegistry.Register("exec", actions.HandlerExec)
	actionRegistry.Register("curl", actions.HandlerCurl)
	actionRegistry.Register("theme", actions.HandlerTheme)
}

// executeEvent ejecuta la accion asociada a un evento declarativo.
// Usa el Registry de acciones para delegar la ejecucion.
func executeEvent(s *AppState, el *ast.ElementNode, eventType string) {
	action, ok := el.Events[eventType]
	if !ok {
		return
	}
	if err := actionRegistry.Execute(action, s.ActionsCB); err != nil {
		log.Printf("Warning: executeEvent(%s): %v", action, err)
	}
}

// executeDataBinding ejecuta bind cuando un input/textarea cambia.
// bind="target-id" actualiza el texto del elemento con ese id.
func executeDataBinding(s *AppState, el *ast.ElementNode) {
	targetID, ok := el.Attrs["bind"]
	if !ok || targetID == "" {
		return
	}
	target := s.FindElementByID(targetID)
	if target == nil {
		return
	}
	target.Text = el.Text
	state.LayoutDirty = true
	s.Dirty = true
}
