package handler

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gdamore/tcell/v2"
)

func TestNew(t *testing.T) {
	e := New()
	if e == nil {
		t.Fatal("New() returned nil")
	}
	if len(e.bindings) != 0 {
		t.Errorf("esperaba bindings vacíos, got %d", len(e.bindings))
	}
	if len(e.actions) != 0 {
		t.Errorf("esperaba actions vacías, got %d", len(e.actions))
	}
}

func TestRegisterAndHandle(t *testing.T) {
	e := New()
	called := false
	e.bindings["d"] = "testAction"
	e.Register("testAction", func() error {
		called = true
		return nil
	})
	ev := tcell.NewEventKey(tcell.KeyRune, 'd', tcell.ModNone)
	if err := e.Handle(ev); err != nil {
		t.Fatalf("Handle() error: %v", err)
	}
	if !called {
		t.Error("la acción no fue ejecutada")
	}
}

func TestHandleUnknownKey(t *testing.T) {
	e := New()
	ev := tcell.NewEventKey(tcell.KeyRune, 'z', tcell.ModNone)
	err := e.Handle(ev)
	if err == nil {
		t.Fatal("esperaba error por tecla sin binding")
	}
}

func TestHandleUnknownAction(t *testing.T) {
	e := New()
	e.bindings["x"] = "nonexistent"
	ev := tcell.NewEventKey(tcell.KeyRune, 'x', tcell.ModNone)
	err := e.Handle(ev)
	if err == nil {
		t.Fatal("esperaba error por acción no registrada")
	}
}

func TestLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "handler.json")
	content := `{
		"version": "1",
		"bindings": [
			{"action": "quit", "keys": ["q", "Ctrl+C"], "description": "Salir"}
		]
	}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}
	e := New()
	if err := e.Load(path); err != nil {
		t.Fatalf("Load() error: %v", err)
	}
	if action, ok := e.bindings["q"]; !ok || action != "quit" {
		t.Errorf("binding 'q' esperado 'quit', got %q", action)
	}
	if action, ok := e.bindings["Ctrl+C"]; !ok || action != "quit" {
		t.Errorf("binding 'Ctrl+C' esperado 'quit', got %q", action)
	}
}

func TestLoad_FileNotFound(t *testing.T) {
	e := New()
	if err := e.Load("/no/existe/handler.json"); err != nil {
		t.Fatal("archivo inexistente no debería ser error:", err)
	}
}
