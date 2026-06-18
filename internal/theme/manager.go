// Gestor de temas con cambio dinamico.
package theme

import (
	"fmt"
	"sync"
)

// Manager gestiona multiples temas y permite cambiar entre ellos.
// Cada tema es un mapa de variables CSS (nombre → valor).
type Manager struct {
	mu     sync.RWMutex
	themes map[string]map[string]string
	Active string // nombre del tema activo
}

// New crea un ThemeManager vacio.
func New() *Manager {
	return &Manager{
		themes: make(map[string]map[string]string),
		Active: "default",
	}
}

// AddTheme agrega o reemplaza un tema.
func (m *Manager) AddTheme(name string, vars map[string]string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.themes[name] = vars
}

// Switch activa un tema por nombre. Retorna error si no existe.
// Devuelve las variables del tema activado.
func (m *Manager) Switch(name string) (map[string]string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	vars, ok := m.themes[name]
	if !ok {
		return nil, fmt.Errorf("tema no encontrado: %s", name)
	}
	m.Active = name
	return vars, nil
}

// GetActiveVars devuelve las variables del tema activo.
func (m *Manager) GetActiveVars() map[string]string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	vars, ok := m.themes[m.Active]
	if !ok {
		return nil
	}
	result := make(map[string]string, len(vars))
	for k, v := range vars {
		result[k] = v
	}
	return result
}

// GetThemeNames devuelve los nombres de temas disponibles.
func (m *Manager) GetThemeNames() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	names := make([]string, 0, len(m.themes))
	for name := range m.themes {
		names = append(names, name)
	}
	return names
}
