// Package handler implementa el engine de bindings de teclado.
package handler

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/gdamore/tcell/v2"
)

// ActionFunc es una función de handler. Retorna error si falla.
type ActionFunc func() error

// Binding asocia teclas con una acción.
type Binding struct {
	Action      string   `json:"action"`
	Keys        []string `json:"keys"`
	Description string   `json:"description"`
}

// Config es la estructura del archivo handler.json.
type Config struct {
	Version  string    `json:"version"`
	Logs     bool      `json:"logs"`
	LogDir   string    `json:"log_dir"`
	AutoCopy bool      `json:"auto_copy"`
	Bindings []Binding `json:"bindings"`
}

// Engine gestiona bindings de teclado y acciones registradas.
type Engine struct {
	bindings map[string]string
	actions  map[string]ActionFunc
	logFile  *os.File
	logMu    sync.Mutex
	autoCopy bool
}

// New crea un Engine vacío.
func New() *Engine {
	return &Engine{
		bindings: make(map[string]string),
		actions:  make(map[string]ActionFunc),
	}
}

// Load carga bindings desde un archivo JSON.
// No es error si el archivo no existe — simplemente no carga bindings.
func (e *Engine) Load(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error leyendo %s: %w", path, err)
	}
	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return fmt.Errorf("error parseando %s: %w", path, err)
	}
	for _, b := range cfg.Bindings {
		for _, key := range b.Keys {
			e.bindings[key] = b.Action
		}
	}
	e.autoCopy = cfg.AutoCopy

	// Configurar logging si está habilitado
	if cfg.Logs {
		logDir := cfg.LogDir
		if logDir == "" {
			logDir = filepath.Join(filepath.Dir(path), "logs")
		}
		if err := os.MkdirAll(logDir, 0755); err != nil {
			return fmt.Errorf("error creando dir de logs %s: %w", logDir, err)
		}
		logPath := filepath.Join(logDir, "handler.log")
		f, err := os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("error abriendo log %s: %w", logPath, err)
		}
		e.logFile = f
	}

	return nil
}

// Register registra una acción por nombre.
func (e *Engine) Register(name string, fn ActionFunc) {
	e.actions[name] = fn
}

// Handle procesa un EventKey tcell y ejecuta la acción correspondiente.
// Retorna error si la tecla no tiene binding o la acción no existe.
func (e *Engine) Handle(ev *tcell.EventKey) error {
	key := FormatKey(ev)
	actionName, ok := e.bindings[key]
	if !ok {
		return fmt.Errorf("no binding for key: %s", key)
	}
	fn, ok := e.actions[actionName]
	if !ok {
		return fmt.Errorf("no action registered: %s", actionName)
	}

	// Loguear la accion si hay log activo
	if e.logFile != nil {
		e.logMu.Lock()
		now := time.Now().Format("2006-01-02 15:04:05")
		fmt.Fprintf(e.logFile, "[%s] %s -> %s\n", now, key, actionName)
		e.logMu.Unlock()
	}

	return fn()
}

// FormatKey convierte un EventKey tcell a string canónico.
// Letras minúsculas, mayúsculas con Shift implícito,
// Ctrl+letra, Alt+letra, y especiales (Esc, Enter, etc.).
func FormatKey(ev *tcell.EventKey) string {
	if ev.Key() != tcell.KeyRune {
		name, ok := keyNames[ev.Key()]
		if ok {
			return name
		}
		return fmt.Sprintf("Key(%d)", ev.Key())
	}
	r := ev.Rune()
	mods := ev.Modifiers()
	if mods&tcell.ModCtrl != 0 {
		return fmt.Sprintf("Ctrl+%c", r)
	}
	if mods&tcell.ModAlt != 0 {
		return fmt.Sprintf("Alt+%c", r)
	}
	return string(r)
}

// AutoCopy retorna si el modo auto-copy esta habilitado.
func (e *Engine) AutoCopy() bool {
	return e.autoCopy
}

// Close cierra el log file si esta abierto.
func (e *Engine) Close() error {
	if e.logFile != nil {
		return e.logFile.Close()
	}
	return nil
}
