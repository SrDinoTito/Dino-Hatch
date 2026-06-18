// Sistema extensible de acciones para Hatch.
// Define tipos base: Callbacks, Handler, Registry.
package actions

import (
	"fmt"
	"strings"
)

// Callbacks agrupa funciones de interaccion con AppState
// para que actions/ no tenga que importar cmd/hatch.
type Callbacks struct {
	Navigate      func(page string)
	ModalOpen     func()
	ModalClose    func()
	ModalToggle   func()
	RandomColors  func()
	Quit          func()
	SetDirty      func()
	ExecOutput    func(command string) (string, error)
	HTTPGet       func(url string) (string, error)
	FindElement    func(id string) interface{}
	SetElementText func(id, text string)
	Log           func(format string, args ...interface{})
	PublishEvent  func(topic, data string)
	SwitchTheme   func(name string) error
}

// Handler procesa una accion tipo "nombre:args".
type Handler func(args string, cb *Callbacks) error

// Registry almacena handlers de acciones registrables.
type Registry struct {
	handlers map[string]Handler
}

// New crea un Registry vacio.
func New() *Registry {
	return &Registry{handlers: make(map[string]Handler)}
}

// Register asocia un nombre con un handler de accion.
func (r *Registry) Register(name string, fn Handler) {
	r.handlers[name] = fn
}

// Execute ejecuta una accion tipo "nombre:args".
// Retorna error si el handler no existe o falla.
func (r *Registry) Execute(action string, cb *Callbacks) error {
	parts := strings.SplitN(action, ":", 2)
	if len(parts) < 2 {
		return fmt.Errorf("formato invalido: %s (esperado: nombre:args)", action)
	}
	cmd, args := parts[0], parts[1]
	fn, ok := r.handlers[cmd]
	if !ok {
		return fmt.Errorf("accion desconocida: %s", cmd)
	}
	return fn(args, cb)
}
