package theme

import (
	"sort"
	"sync"
	"testing"
)

// 1. New() crea manager con Active="default" y themes vacío.
func TestNew(t *testing.T) {
	m := New()
	if m == nil {
		t.Fatal("New() returned nil")
	}
	if m.Active != "default" {
		t.Errorf("Active = %q, want %q", m.Active, "default")
	}
	if m.themes == nil {
		t.Fatal("themes map is nil")
	}
	if len(m.themes) != 0 {
		t.Errorf("len(themes) = %d, want 0", len(m.themes))
	}
}

// 6. Switch cambia Active correctamente.
// 2. AddTheme + GetActiveVars (agregar tema, switchear, verificar vars).
func TestAddThemeAndSwitch(t *testing.T) {
	m := New()
	vars := map[string]string{"--bg": "#000", "--fg": "#fff"}
	m.AddTheme("dark", vars)

	// Switch activa el tema y devuelve las variables
	result, err := m.Switch("dark")
	if err != nil {
		t.Fatalf("Switch('dark') error: %v", err)
	}
	if m.Active != "dark" {
		t.Errorf("Active = %q, want %q", m.Active, "dark")
	}
	if result["--bg"] != "#000" || result["--fg"] != "#fff" {
		t.Errorf("Switch returned vars = %v, want %v", result, vars)
	}

	// GetActiveVars coincide
	active := m.GetActiveVars()
	if active["--bg"] != "#000" || active["--fg"] != "#fff" {
		t.Errorf("GetActiveVars() = %v, want %v", active, vars)
	}

	// Agregar otro tema y switchear de nuevo
	vars2 := map[string]string{"--bg": "#fff", "--fg": "#000"}
	m.AddTheme("light", vars2)
	r2, err := m.Switch("light")
	if err != nil {
		t.Fatalf("Switch('light') error: %v", err)
	}
	if m.Active != "light" {
		t.Errorf("Active = %q, want %q", m.Active, "light")
	}
	if r2["--bg"] != "#fff" {
		t.Errorf("Switch returned --bg = %q, want %q", r2["--bg"], "#fff")
	}
	active2 := m.GetActiveVars()
	if active2["--bg"] != "#fff" || active2["--fg"] != "#000" {
		t.Errorf("GetActiveVars after switch = %v, want %v", active2, vars2)
	}
}

// 3. Switch a tema inexistente → error.
func TestSwitchNonExistent(t *testing.T) {
	m := New()
	_, err := m.Switch("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent theme, got nil")
	}
	// Active no debe haber cambiado
	if m.Active != "default" {
		t.Errorf("Active changed to %q after failed switch", m.Active)
	}
}

// 4. GetActiveVars devuelve copia (modificar resultado no afecta original).
func TestGetActiveVarsReturnsCopy(t *testing.T) {
	m := New()
	vars := map[string]string{"--bg": "red", "--fg": "blue"}
	m.AddTheme("mytheme", vars)
	m.Switch("mytheme")

	got := m.GetActiveVars()
	// Modificar la copia
	got["--bg"] = "green"
	got["--new"] = "added"

	// El original no debe haber cambiado
	original := m.GetActiveVars()
	if original["--bg"] != "red" {
		t.Errorf("original --bg = %q, want %q (copy was mutated)", original["--bg"], "red")
	}
	if _, exists := original["--new"]; exists {
		t.Error("original was mutated: --new key appeared")
	}
}

// 5. GetThemeNames (agregar 2 temas, verificar nombres).
func TestGetThemeNames(t *testing.T) {
	m := New()
	// Solo "default" en el map, no como tema registrado
	names := m.GetThemeNames()
	if len(names) != 0 {
		t.Errorf("GetThemeNames() = %v, want empty slice", names)
	}

	m.AddTheme("dark", map[string]string{"--a": "1"})
	m.AddTheme("light", map[string]string{"--b": "2"})

	names = m.GetThemeNames()
	sort.Strings(names)
	expected := []string{"dark", "light"}
	if len(names) != len(expected) {
		t.Fatalf("GetThemeNames() = %v, want %v", names, expected)
	}
	for i, name := range names {
		if name != expected[i] {
			t.Errorf("GetThemeNames()[%d] = %q, want %q", i, name, expected[i])
		}
	}
}

// 7. GetActiveVars cuando Active="default" pero no existe en el map → nil.
func TestGetActiveVarsNilWhenNoTheme(t *testing.T) {
	m := New()
	// Active es "default" pero nunca se agregó un tema "default"
	vars := m.GetActiveVars()
	if vars != nil {
		t.Errorf("GetActiveVars() = %v, want nil", vars)
	}

	// También después de Switchear a un tema que luego se borra (simulado)
	m.AddTheme("temp", map[string]string{"--x": "y"})
	m.Switch("temp")
	// No podemos borrar temas, pero podemos hacer Switch a uno que no existe
	// para probar el path: Active apunta a tema existente → funciona
	vars = m.GetActiveVars()
	if vars == nil {
		t.Error("GetActiveVars() returned nil after Switch to existing theme")
	}
}

// 8. Thread-safety: goroutines concurrentes agregando mientras otro lee.
func TestThreadSafety(t *testing.T) {
	m := New()
	var wg sync.WaitGroup

	// Escritores concurrentes
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			name := "dark"
			if i%2 == 0 {
				name = "light"
			}
			m.AddTheme(name, map[string]string{
				"--bg": "value",
			})
		}(i)
	}

	// Lectores concurrentes
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			m.GetActiveVars()
		}()
	}

	// Switch concurrente
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			names := []string{"dark", "light", "nonexistent"}
			name := names[i%len(names)]
			m.Switch(name) // ignoro error
		}(i)
	}

	wg.Wait()

	// GetThemeNames concurrente
	var wg2 sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			m.GetThemeNames()
		}()
	}
	wg2.Wait()
}

// TestSwitchReturnsSameMapAsActiveVars: Switch devuelve las vars y
// GetActiveVars devuelve una copia consistente (verifica igualdad).
func TestSwitchReturnsConsistentWithActiveVars(t *testing.T) {
	m := New()
	expected := map[string]string{"--color": "blue", "--size": "12"}
	m.AddTheme("test", expected)
	gotSwitch, err := m.Switch("test")
	if err != nil {
		t.Fatalf("Switch error: %v", err)
	}
	gotActive := m.GetActiveVars()

	// Switch devuelve el mapa interno (sin copia), GetActiveVars una copia
	for k, v := range expected {
		if gotSwitch[k] != v {
			t.Errorf("Switch[%q] = %q, want %q", k, gotSwitch[k], v)
		}
		if gotActive[k] != v {
			t.Errorf("GetActiveVars[%q] = %q, want %q", k, gotActive[k], v)
		}
	}
}

// TestAddThemeReplacesExisting: AddTheme reemplaza vars de un tema existente.
func TestAddThemeReplacesExisting(t *testing.T) {
	m := New()
	m.AddTheme("x", map[string]string{"--old": "1"})
	m.Switch("x")

	m.AddTheme("x", map[string]string{"--new": "2"})
	// Switch otra vez para obtener las nuevas vars
	vars, err := m.Switch("x")
	if err != nil {
		t.Fatalf("Switch error: %v", err)
	}
	if _, exists := vars["--old"]; exists {
		t.Error("old var still present after theme replacement")
	}
	if vars["--new"] != "2" {
		t.Errorf("--new = %q, want %q", vars["--new"], "2")
	}
}

// TestGetActiveVarsMultipleThemes: múltiples temas, GetActiveVars refleja el activo.
func TestGetActiveVarsMultipleThemes(t *testing.T) {
	m := New()
	dark := map[string]string{"--bg": "black"}
	light := map[string]string{"--bg": "white"}
	m.AddTheme("dark", dark)
	m.AddTheme("light", light)

	m.Switch("dark")
	if v := m.GetActiveVars()["--bg"]; v != "black" {
		t.Errorf("active bg = %q, want %q", v, "black")
	}

	m.Switch("light")
	if v := m.GetActiveVars()["--bg"]; v != "white" {
		t.Errorf("active bg = %q, want %q", v, "white")
	}

	// No debe mezclarse
	if len(m.GetActiveVars()) != 1 {
		t.Errorf("expected 1 var, got %d", len(m.GetActiveVars()))
	}
}

// TestGetThemeNamesDoesNotIncludeActive: GetThemeNames solo devuelve temas
// agregados con AddTheme, no depende de Active.
func TestGetThemeNamesOnlyAddedThemes(t *testing.T) {
	m := New()
	// Active="default", pero nunca agregado
	names := m.GetThemeNames()
	for _, n := range names {
		if n == "default" {
			t.Error("GetThemeNames should not include 'default' before AddTheme")
		}
	}

	m.AddTheme("default", map[string]string{"--a": "1"})
	names = m.GetThemeNames()
	found := false
	for _, n := range names {
		if n == "default" {
			found = true
		}
	}
	if !found {
		t.Error("GetThemeNames should include 'default' after AddTheme")
	}
}
