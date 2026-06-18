# HATCH-CORE — Design

## Version
0.1.0-alpha

## Referencias

| Referencia | Descripción |
|-----------|-------------|
| REQ-001 | Parseo de HML |
| REQ-002 | Parseo de HSS |
| REQ-003 | AST con estilos computados |
| REQ-004 | Layout flexbox |
| REQ-005 | Renderizado a terminal |
| REQ-006 | Cell buffer con diff persistente |
| REQ-007 | Línea de comandos |
| REQ-008 | Dependencias mínimas |
| REQ-009 | Comentarios en español |
| REQ-013 | Component system (include) |
| REQ-014 | Overflow scroll en contenedores |
| REQ-015 | Eventos declarativos |
| REQ-016 | Data binding |
| REQ-017 | CSS variables / theming |
| REQ-018 | Tab focus navigation |

---

## Arquitectura General

```
┌─────────────┐     ┌──────────────┐     ┌─────────────────────┐
│  .hml file  │ ──→ │   PARSER     │ ──→ │  AST raw + reglas   │
│  (HML+HSS)  │     │  (hml+hss)   │     │  HSS + CSS vars     │
└─────────────┘     └──────────────┘     └──────────┬──────────┘
                                                     │
                                                     ▼
┌───────────────────┐     ┌────────────────┐     ┌─────────────────┐
│  RENDER (tcell)   │ ←── │  LAYOUT ENGINE │ ←── │  COMPUTE STYLES │
│  - CellBuffer     │     │  (flexbox +    │     │  - resolveInclude│
│  - Diff()         │     │   clamp +      │     │  - CSS vars      │
│  - scrollbar      │     │   scroll)      │     │  - herencia      │
│  - clipping       │     └────────────────┘     └─────────────────┘
└───────────────────┘
        │
        ▼
┌──────────────┐
│   Terminal   │
│  (tcell v3)  │
└──────────────┘
```

### Pipeline extendido
1. **Read**: Leer archivo .hml del disco
2. **Parse HML**: Parser produce AST raw con nodos `<include>` sin resolver
3. **Parse HSS**: Bloques `<style>` → StyleRules + ParseCSSVars extrae :root vars
4. **ComputeStyles**: Merge estilos + resolución CSS vars + herencia
5. **LoadPage** (cmd/hatch):
   a. Parsear página .hml
   b. **resolveIncludes**: Reemplazar `<include>` por AST del componente
   c. Merge estilos del componente en reglas globales
   d. Merge CSS vars de la página
   e. ComputeStyles con vars globales
   f. **buildFocusOrder**: Construir orden de Tab
6. **Layout**: Calcular BoundBox (x, y, w, h) para cada nodo
7. **Render**: CellBuffer persistente → Diff() → tcell.SetContent()
8. **Event Loop**: PollEvent → dispatch (mouse/key/resize) → dirty flag → goto 6

---

## DSG-001 — AST (REQ-001, REQ-002, REQ-003)

### Tipos principales

```go
// Document representa el archivo .hml completo
type Document struct {
    Pages     []Page
    ThemeVars map[string]string // CSS variables de :root, ej: {"--bg": "#333"}
}

// Page representa una pantalla/interfaz
type Page struct {
    Name     string
    Width    int          // 0 = auto/terminal width
    Height   int          // 0 = auto/terminal height
    Style    ComputedStyle
    Children []ElementNode
}

// ElementNode representa un elemento UI
type ElementNode struct {
    Tag        string            // "box", "text", "input", "textarea", "button", "include"
    Attrs      map[string]string // atributos raw del XML (sin eventos, sin id)
    Style      ComputedStyle
    Children   []ElementNode
    Text       string            // contenido textual para <text>, <button>, <input>, <textarea>
    BoundBox   BoundBox          // calculado en fase layout
    ID         string            // atributo id
    ScrollX    int               // scroll offset horizontal (para overflow:scroll)
    ScrollY    int               // scroll offset vertical (para overflow:scroll)
    Events     map[string]string // "click" -> "page:proyectos", "change" -> "..."
    IncludeSrc string            // "" si no es include, "ruta/archivo.hml" si es include
}

// ComputedStyle con todas las propiedades resueltas.
type ComputedStyle struct {
    Width, Height       int
    Grow                float64
    Gap                 int
    Align, Justify      string  // "start", "center", "end", "stretch"
    Direction           string  // "row", "column"
    Color, BgColor      tcell.Color
    Border              bool
    TextAlign           string  // "left", "center", "right"
    VAlign              string  // "top", "middle", "bottom"
    Padding, Margin     int
    MinWidth, MinHeight int     // 0 = sin constraint
    MaxWidth, MaxHeight int     // 0 = sin constraint
    Overflow            string  // "visible" | "hidden" | "scroll" (default "visible")
}

// BoundBox posicion calculada por layout
type BoundBox struct {
    X, Y, W, H int
}

// StyleRule representa una regla HSS sin resolver
type StyleRule struct {
    Selector   string
    Properties map[string]string
}
```

### Decisiones
- **encoding/xml** de stdlib para parsear HML (REQ-008). Se preprocesa el contenido
  eliminando bloques `<style>` antes del parseo XML (`preprocessStyleBlocks` en hml.go).
- HSS se parsea manualmente con un scanner simple (split + trim), no hay librería Go
  que haga CSS parsing ligero.
- Los estilos computados se resuelven en 3 pasos: (1) defaults, (2) reglas HSS por
  selector, (3) atributos inline (Attrs). Herencia del padre para propiedades no
  definidas.
- Atributos de eventos (onclick, onchange, onfocus, onblur) se extraen de Attrs a
  Events durante el parseo en `parseElement()`.
- El tag `<include>` se parsea como un nodo normal con `IncludeSrc` no vacío; la
  resolución ocurre post-parseo en `resolveIncludes()`.

---

## DSG-002 — Layout Flexbox Simplificado (REQ-004)

### Algoritmo
1. El contenedor tiene un ancho/alto conocido (terminal size o calculado)
2. Se restan border (2px) y padding (pad*2) del área disponible
3. Se suman los grow de los hijos y los tamaños fijos
4. El espacio restante (total - fixed - gaps - margins) se distribuye proporcionalmente
   al grow
5. gap se aplica entre hijos; margin se aplica alrededor de cada hijo
6. align/justify se aplican según el eje transversal/principal
7. Clamping min/max: después de asignar tamaños, se ajustan a min/max-width/height
8. Post-expansion: si un contenedor no tiene width/height explícito, se expande
   para cubrir a sus hijos (con estabilización selectiva para textarea)

### Flexbox properties soportadas
| Propiedad | Valores | Default |
|-----------|---------|---------|
| direction | row, column | column |
| grow | float | 0 |
| gap | int | 0 |
| align | start, center, end, stretch | stretch |
| justify | start, center, end, space-between | start |
| width, height | int | 0 (auto) |
| padding, margin | int | 0 |
| min-width, min-height | int | 0 (sin constraint) |
| max-width, max-height | int | 0 (sin constraint) |
| overflow | visible, hidden, scroll | visible |

### ContentHeight / ContentWidth
- `ContentHeight(n)` calcula altura total del contenido intrínseco de un elemento
  (suma alturas de hijos + gaps + padding). Usado para scroll containers y MaxScroll.
- `ContentWidth(n)` equivalente para ancho.
- Cache: solo se calcula en dirty frames (B4).

---

## DSG-003 — Cell Buffer y Diff (REQ-005, REQ-006)

```
Frame anterior (PrevCB) ─┐
                          ├──→ Diff() → []Update → tcell.SetContent()
Frame actual  (CurrCB)  ──┘
```

- **CellBuffer**: struct con grid 2D de Cell{rune, style}. Métodos: Set, Get, Fill,
  Resize, Diff.
- **CellBuffer persistente**: AppState tiene PrevCB y CurrCB. Se mantienen entre
  frames y se intercambian (swap) después de cada render.
- **Diff()**: Compara dos CellBuffer celda por celda. Retorna solo las celdas que
  cambiaron ([]Update{ X, Y, Cell }). Permite enviar solo cambios a tcell.
- **Dirty flag (B2)**: AppState.Dirty = true desencadena layout/render completo.
  Si Dirty=false, el event loop salta el render y solo espera eventos.
- **Mouse throttle (B3)**: Para mouse-move, solo se marca Dirty si pasaron >=33ms
  desde el último evento (≈30 FPS idle).
- **Resize**: En EventResize, se crean nuevos buffers con el nuevo tamaño y se
  fuerza Dirty=true.

---

## DSG-004 — tcell Screen Driver (REQ-005)

- Se implementa la interfaz `Screen` en `internal/render/screen.go`:
  - `Init()`, `Flush()`, `Close()`, `Size()` (width, height int, err)
- `tcellScreen`: implementa Screen usando `tcell.NewScreen()` directamente.
- `mockScreen`: implementa Screen para tests, almacena celdas en buffer interno.
- El event loop principal en `cmd/hatch/eventloop.go` usa `tcell.Screen` directamente
  (no la interfaz abstracta) para tener acceso a PollEvent, SetContent, Show.

---

## DSG-005 — Component System (REQ-013)

### Flujo de resolución de includes

```
LoadPage(name):
  1. Parsear pagina.hml → pageDoc (AST raw con nodos <include>)
  2. resolveIncludes(pageDoc):
     a. Recorrer AST con walkAndResolve()
     b. Para cada nodo con IncludeSrc != "":
        i.  Leer archivo del componente (ruta relativa al .hml actual)
        ii. Parsear componente → compDoc
        iii. Extraer styles del componente → mergear en s.StyleRules
        iv. Mergear atributos (los del include ganan sobre los del componente)
        v.  Si el include tiene id, asignarlo al root del componente
        vi. Resolver includes anidados dentro del componente
        vii. Reemplazar nodo <include> por los children del componente
  3. Mergear CSS vars de la pagina en s.Doc.ThemeVars
  4. Mergear StyleRules de la pagina con las globales
  5. ComputeStyles(s.Doc, mergedRules, s.Doc.ThemeVars)
  6. buildFocusOrder()
```

### mergeAttrs
```go
func mergeAttrs(target, source map[string]string) {
    for k, v := range source { target[k] = v } // source gana
}
```

### Decisiones
- La resolución ocurre después del parseo HML y antes de ComputeStyles, porque
  los estilos del componente deben estar disponibles para el cálculo de estilos.
- Includes anidados se resuelven recursivamente con `walkAndResolve`.
- Si un archivo include no se encuentra, se emite warning y se omite (no crash).
- Los atributos del tag `<include>` se mergean con los del root del componente,
  permitiendo sobrescribir id, estilos inline, etc.

---

## DSG-006 — Overflow Scroll Containers (REQ-014)

### Scroll interno por contenedor
- Cada `ElementNode` tiene `ScrollY int` para scroll vertical independiente.
- El scroll de página global es `AppState.ScrollY`.
- La propiedad `overflow` ("scroll" | "hidden" | "visible") determina el
  comportamiento.

### Clipping de hijos
- `isChildVisible(child, parentScrollY, parentHeight, parentY) bool`: verifica
  si un hijo está dentro del viewport del padre, considerando el scroll offset.
- `clipHeight(child, parentScrollY, parentHeight, parentY) int`: calcula la
  altura visible de un hijo, limitando top/bottom al viewport del padre.
- En `render.go`, los hijos de contenedores con overflow=scroll/hidden se
  renderizan con `isChildVisible` para saltar hijos fuera del viewport.

### Scrollbar visual en contenedor
- `drawContainerScrollbar(cb, el, pageScrollY)`: dibuja barra proporcional
  en el margen derecho del contenedor.
- `drawAllContainerScrollbars(cb, doc, pageScrollY)`: recorre AST y dibuja
  scrollbars para todos los contenedores con overflow=scroll.

### Wheel event routing
- En `mouse_handler.go`, wheel events primero buscan un contenedor scrollable
  bajo el cursor con `findScrollContainer()`. Si se encuentra, se scrollea
  internamente (ScrollY += 3). Si no, se scrollea la página global.

### findScrollContainer
- `findScrollContainer(page, mx, my, pageScrollY)`: DFS que encuentra el
  contenedor overflow=scroll más profundo bajo las coordenadas del mouse.

---

## DSG-007 — Declarative Events + Data Binding (REQ-015, REQ-016)

### Eventos declarativos
Los atributos onclick, onchange, onfocus, onblur se extraen de Attrs a Events
durante el parseo en `hml_parse.go`. El formato es `tipo:argumento`.

**Acciones soportadas** (en `handler.go`):
| Formato | Ejemplo | Comportamiento |
|---------|---------|---------------|
| `page:NOMBRE` | `onclick="page:proyectos"` | Navega a la página (LoadPage) |
| `modal:open` | `onclick="modal:open"` | Abre modal overlay |
| `modal:close` | `onclick="modal:close"` | Cierra modal |
| `modal:toggle` | `onclick="modal:toggle"` | Alterna modal |
| `action:quit` | `onclick="action:quit"` | Sale de la aplicación |
| `action:random_colors` | `onclick="action:random_colors"` | Activa/desactiva colores aleatorios |

**Flujo de ejecución**:
1. `handleMousePress` detecta click en elemento con `Events["click"] != ""`
2. Llama a `executeEvent(s, el, "click")`
3. `executeEvent` parsea `tipo:argumento` y ejecuta la acción correspondiente

### Data binding
El atributo `bind="target-id"` en elementos input/textarea sincroniza el texto
del elemento con el target especificado.

**Puntos de sincronización** (en `input.go`):
- Tecla Rune insertada → `executeDataBinding(state, el)`
- Backspace → `executeDataBinding(state, el)`
- Delete → `executeDataBinding(state, el)`
- Enter en textarea → `executeDataBinding(state, el)`

**Flujo**:
```go
func executeDataBinding(s *AppState, el *ast.ElementNode) {
    targetID := el.Attrs["bind"]
    target := s.FindElementByID(targetID)
    target.Text = el.Text  // sincroniza texto
    s.Dirty = true          // fuerza re-render
}
```

---

## DSG-008 — CSS Variables / Theming (REQ-017)

### Extracción de variables
`ParseCSSVars(styleContent string) map[string]string` escanea bloques `:root { ... }`
y extrae pares `--nombre: valor`. Se llama en dos lugares:
- `main.go`: al parsear el archivo .hml inicial
- `pipeline.go` (LoadPage): al cargar cada página, mergeando vars nuevas en
  `s.Doc.ThemeVars`

### Resolución de variables
`resolveCSSVars(val string, vars map[string]string) string` reemplaza ocurrencias
de `var(--nombre)` con el valor del mapa. Si la variable no existe, se deja
sin reemplazar (el parser ignorará la propiedad).

### Propagación
Las variables se pasan a `ComputeStyles(doc, rules, vars)` y luego a
`applyProps(s, props, vars, exp)` que resuelve `var()` en cada valor antes de
aplicarlo al ComputedStyle.

### Flujo completo
```go
// En main.go o LoadPage:
vars := parser.ParseCSSVars(allStyleContent)
doc.ThemeVars = vars
doc = parser.ComputeStyles(doc, rules, doc.ThemeVars)
// → resolveNode → applyProps → resolveCSSVars por cada propiedad
```

---

## DSG-009 — Tab Focus Navigation (REQ-018)

### FocusOrder
`buildFocusOrder()` recorre el AST en DFS y recolecta elementos focusables
(input, textarea, button) en `AppState.FocusOrder []*ast.ElementNode`.

```go
func (s *AppState) buildFocusOrder() {
    s.FocusOrder = nil
    var walk func(n *ast.ElementNode)
    walk = func(n *ast.ElementNode) {
        if n.Tag == "input" || n.Tag == "textarea" || n.Tag == "button" {
            s.FocusOrder = append(s.FocusOrder, n)
        }
        for i := range n.Children { walk(&n.Children[i]) }
    }
    // Recorrer todas las páginas
}
```

### FocusIndex
`AppState.FocusIndex int` (-1 = ninguno) mantiene la posición actual en el orden.

### Navegación Tab/Shift+Tab (en keyboard.go)
```go
case tcell.KeyTab:
    if e.Modifiers()&tcell.ModShift != 0 {
        s.FocusIndex-- // Shift+Tab: reversa
    } else {
        s.FocusIndex++ // Tab: avance
    }
    // Cíclico: si sale del rango, vuelve al extremo opuesto
    s.FocusedElement = s.FocusOrder[s.FocusIndex]
```

### Comportamiento
- Al recibir foco, input/textarea muestran cursor (`▌`)
- El cursor se posiciona al final del texto existente (o en la posición del click)
- keyboard.go maneja teclas de navegación (arrows, pgup/pgdn, home/end) para
  scroll global de página
- Si hay un modal abierto, Tab/Shift+Tab no navega (return inmediato)
- Enter en textarea se maneja antes que el dispatch general de teclas

---

## Integración

### Diagrama de paquetes completo
```
cmd/hatch/                          ← Entry point CLI (22 archivos)
  main.go                           ← Flags → parse → tcell → RunLoop
  pipeline.go                       ← AppState, LoadPage, NewAppState, resolveIncludes
  eventloop.go                      ← RunLoop (dirty flag, diff, scrollbars)
  render.go                         ← renderDoc, renderPage, renderNode
  render_elements.go                ← renderBoxContent, renderButton, renderText, etc.
  render_helpers.go                 ← drawBorder, drawOverlay, renderTextarea
  render_clip.go                    ← isChildVisible, clipHeight
  scrollbar.go                      ← drawScrollbar, drawContainerScrollbar
  scroll_container.go               ← findScrollContainer
  textarea_scrollbar.go             ← renderTextareaScrollbar
  handler.go                        ← executeEvent, executeDataBinding
  include.go                        ← resolveIncludes, walkAndResolve, mergeAttrs
  navigate.go                       ← ContentHeight, FindElementByID
  focus.go                          ← buildFocusOrder
  input.go                          ← handleInputKey, cursorLineCol, moveCursorUp/Down
  interactive.go                    ← hitTest, copyToClipboard, getElementText
  mouse_handler.go                  ← handleMouseEvent, handleMousePress
  mouse_release.go                  ← handleMouseRelease
  keyboard.go                       ← handleKeyEvent (Tab, arrows, scroll)
  colors.go                         ← toggleRandomColors, collectBoxes
  state.go                          ← inputState, normalizedSelectionRect
  diag_test.go                      ← Tests de textarea/layout

internal/
  ast/
    node.go                         ← Tipos AST (Document, Page, ElementNode, ...)
    node_test.go                    ← Tests AST
  parser/
    hml.go                          ← ParseHML entry point, preprocessStyleBlocks
    hml_parse.go                    ← pageFromStartElement, parseElement, attrsFromSlice
    hss.go                          ← ParseHSS, ParseCSSVars, parseProperties, isKnownProp
    compute.go                      ← ComputeStyles, resolveNode, parseInt, parseColor, etc.
    props.go                        ← applyProps, resolveCSSVars
    inherit.go                      ← inheritProps
    *._test.go                      ← Tests (hml, hss, compute, props, inherit)
    testdata/                       ← Fixtures .hml
  handler/
    engine.go                       ← Engine (Load, Register, Handle, FormatKey)
    keys.go                         ← keyNames map
    engine_test.go, keys_test.go, logging_test.go
  layout/
    flex.go                         ← Layout, layoutPage, layoutNode, post-expansion
    layout.go                       ← layoutChildren (algoritmo flexbox principal)
    content.go                      ← ContentHeight, ContentWidth
    measure.go                      ← intrinsicSize
    *._test.go                      ← Tests (flex, basic, distribution, align, edge, helpers)
  render/
    cellbuffer.go                   ← CellBuffer, Diff, Resize, Fill
    cellbuffer_test.go              ← Tests cell buffer
    screen.go                       ← Screen interface + tcell driver
    screen_test.go                  ← Tests con mockScreen

canva/                              ← Archivos de ejemplo/demo
  demo.hml                          ← HML principal: layout base + navbar + sidebar
  handler.json                      ← Bindings de teclado
  components/                       ← 9 componentes reutilizables
    header.hml, modal.hml, scroll_demo.hml, events_demo.hml,
    bind_demo.hml, theme_demo.hml, textarea_demo.hml,
    layout_demo.hml, tabs_demo.hml
  pages/                            ← 4 páginas de la app demo
    inicio.hml (dashboard), proyectos.hml (layout+scroll),
    config.hml (bind+events+theme), ayuda.hml (textarea+tabs+atajos)
```

### Flujo de inicio completo
```
main.go:
  1. Parse flags (hatch run <file>)
  2. Read file
  3. Parse HML → doc, styleBlocks
  4. Parse HSS → rules
  5. ParseCSSVars → doc.ThemeVars
  6. ComputeStyles(doc, rules, vars)
  7. Initialize tcell screen
  8. NewAppState(doc, rules, scr):
     a. Cargar modal.hml como ModalDoc
     b. Registrar acciones globales (quit, random_colors)
     c. Cargar handler.json (bindings de teclado)
     d. LoadPage("inicio"):
        - Parsear pagina .hml
        - resolveIncludes (reemplazar <include>)
        - Merge CSS vars + StyleRules
        - ComputeStyles
        - buildFocusOrder
     e. Crear CurrCB y PrevCB
     f. Dirty = true
  9. RunLoop(state):
     while Running:
       if Dirty:
         a. Layout(doc, W, H)
         b. Calcular MaxScroll = ContentHeight() - H
         c. Render a CurrCB (con clipping y scrollbars)
         d. Diff() PrevCB vs CurrCB → SetContent()
         e. Swap buffers, Show(), Dirty = false
       PollEvent:
         mouse  → handleMouseEvent (throttle)
         resize → resize buffers, Dirty = true
         key    → handleKeyEvent, Dirty = true
```

---

## Historial de cambios
| Version | Fecha | Cambios |
|---------|-------|---------|
| 0.0.1-draft | - | Versión inicial: parseo, layout básico, render |
| 0.1.0-alpha | Jun 2026 | Padding/margin/minmax, refactor, cellbuf persistente, dirty flag, scrollbar, component system, overflow scroll, eventos, data binding, CSS vars, Tab focus, demo rediseñada |
