// Handler para cambio de temas via ThemeManager.
package actions

// HandlerTheme cambia el tema activo via Callbacks.SwitchTheme.
// Formato: "theme:dark" / "theme:light" / "theme:toggle"
func HandlerTheme(args string, cb *Callbacks) error {
	switch args {
	case "toggle":
		// toggle entre el primero y segundo tema disponible
		cb.Log("theme: toggle no implementado aun sin lista de temas")
		return nil
	default:
		if cb.SwitchTheme != nil {
			return cb.SwitchTheme(args)
		}
	}
	return nil
}
