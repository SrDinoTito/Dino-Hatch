package handler

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gdamore/tcell/v2"
)

// TestLogging registra debug_borders, settea logs:true, ejecuta Handle y verifica que se escribio el log.
func TestLogging(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "handler.json")
	jsonData := `{
		"version": "1",
		"logs": true,
		"bindings": [
			{"action": "debug_borders", "keys": ["D"], "description": ""}
		]
	}`
	if err := os.WriteFile(jsonPath, []byte(jsonData), 0644); err != nil {
		t.Fatal(err)
	}

	eng := New()
	toggled := false
	eng.Register("debug_borders", func() error {
		toggled = !toggled
		return nil
	})
	if err := eng.Load(jsonPath); err != nil {
		t.Fatal(err)
	}
	defer eng.Close()

	// Simular Shift+D
	ev := tcell.NewEventKey(tcell.KeyRune, 'D', tcell.ModNone)
	if err := eng.Handle(ev); err != nil {
		t.Fatal(err)
	}

	// Verificar que se creo el log
	logContent, err := os.ReadFile(filepath.Join(tmpDir, "logs", "handler.log"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(logContent), "D -> debug_borders") {
		t.Fatalf("log no contiene la entrada esperada: %s", logContent)
	}
}

// TestLoggingOff verifica que con logs:false no se crea archivo de log.
func TestLoggingOff(t *testing.T) {
	tmpDir := t.TempDir()
	jsonPath := filepath.Join(tmpDir, "handler.json")
	jsonData := `{
		"version": "1",
		"logs": false,
		"bindings": [
			{"action": "debug_borders", "keys": ["D"], "description": ""}
		]
	}`
	if err := os.WriteFile(jsonPath, []byte(jsonData), 0644); err != nil {
		t.Fatal(err)
	}

	eng := New()
	eng.Register("debug_borders", func() error { return nil })
	if err := eng.Load(jsonPath); err != nil {
		t.Fatal(err)
	}
	defer eng.Close()

	ev := tcell.NewEventKey(tcell.KeyRune, 'D', tcell.ModNone)
	_ = eng.Handle(ev)

	// Verificar que NO se creo el dir logs
	if _, err := os.Stat(filepath.Join(tmpDir, "logs")); !os.IsNotExist(err) {
		t.Fatal("logs/ se creo a pesar de logs:false")
	}
}
