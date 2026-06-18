# Glosario — dino-hatch v0.1.0

> Referencia completa de elementos HML, propiedades HSS, eventos, acciones y APIs del framework.

## HML Tags

| Tag | Descripción | Hijos | Atributos clave |
|-----|-------------|-------|-----------------|
| `<page>` | Página raíz del documento. Todo `.hml` debe tener al menos una. | `box`, `text`, `button`, `input`, `textarea`, `include` | `name`, `width`, `height` |
| `<box>` | Contenedor genérico con layout flexbox. | Cualquier tag | `id`, `style`, `onclick`, `onchange`, `onfocus`, `onblur` |
| `<text>` | Texto plano (no puede tener hijos). | Ninguno | `id`, `style` |
| `<button>` | Botón interactivo. | Ninguno (texto inline) | `id`, `style`, `onclick` |
| `<input>` | Campo de texto de una línea. | Ninguno | `id`, `style`, `bind`, `onchange`, `onfocus`, `onblur` |
| `<textarea>` | Área de texto multilínea con scroll interno. | Ninguno | `id`, `style`, `bind`, `onchange`, `onfocus`, `onblur` |
| `<include>` | Inclusión de componente reutilizable. | Consumido (no se parsea) | `src` (ruta al .hml) |

## Atributos Globales

| Atributo | Aplica a | Descripción |
|----------|----------|-------------|
| `id` | Todos | Identificador único para data binding y navegación |
| `style` | Todos | Estilos HSS inline (ej: `style="color: red; padding: 2"`) |
| `class` | Todos | Clase CSS para estilos compartidos |
| `onclick` | box, button, text, input, textarea | Evento al hacer clic |
| `onchange` | input, textarea | Evento al cambiar el valor |
| `onfocus` | input, textarea | Evento al recibir foco |
| `onblur` | input, textarea | Evento al perder foco |
| `bind` | input, textarea | Sincroniza el valor con el texto de otro elemento por ID |
| `src` | include | Ruta al archivo .hml del componente |

## HSS Properties (Bloques `<style>` o inline)

### Dimensiones y Espaciado

| Propiedad | Valores | Descripción |
|-----------|---------|-------------|
| `width` | número (caracteres) | Ancho fijo del elemento |
| `height` | número (líneas) | Alto fijo del elemento |
| `min-width` | número | Ancho mínimo (clamping) |
| `min-height` | número | Alto mínimo (clamping) |
| `max-width` | número | Ancho máximo (clamping) |
| `max-height` | número | Alto máximo (clamping) |
| `padding` | número | Espaciado interno (mismos px en 4 lados) |
| `margin` | número | Espaciado externo entre hermanos (eje primario) |
| `border` | `true` / `false` | Borde de 1 caracter alrededor del contenido |

### Layout Flexbox

| Propiedad | Valores | Descripción |
|-----------|---------|-------------|
| `direction` | `row` / `column` | Eje principal del layout |
| `grow` | número (0=no crece) | Factor de crecimiento para espacio extra |
| `gap` | número | Espacio entre hijos |
| `align` | `stretch` / `start` / `center` / `end` | Alineación en eje transversal |
| `justify` | `start` / `center` / `end` / `space-between` | Alineación en eje principal |
| `overflow` | `visible` / `scroll` | Comportamiento de desbordamiento |

### Colores y Estilo

| Propiedad | Valores | Descripción |
|-----------|---------|-------------|
| `color` | nombre / `#rrggbb` / `var(--name)` | Color de texto |
| `background-color` | nombre / `#rrggbb` / `var(--name)` | Color de fondo |
| `font-weight` | `bold` / `normal` | Peso de fuente (negrita) |

## CSS Variables / Theming

```hml
<style>
:root {
  --primary: #00ff00;
  --bg: #1a1a2e;
}
</style>
```

| Sintaxis | Descripción |
|----------|-------------|
| `:root { --name: value }` | Define una variable CSS en el tema |
| `var(--name)` | Usa el valor de la variable en cualquier propiedad |
| `theme.add("name", vars)` | API Go para registrar un tema programáticamente |
| `theme.switch("name")` | Cambia el tema activo en runtime |

## Acciones (Valores de onclick)

Las acciones usan formato `tipo:argumento`.

| Acción | Formato | Descripción |
|--------|---------|-------------|
| Navegación | `page:nombre` | Cambia a otra página del documento |
| Modal | `modal:open` / `modal:close` | Abre/cierra el modal overlay |
| Acción registrada | `action:nombre(args)` | Ejecuta un handler registrado via `Registry` |
| Comando shell | `exec:comando --args` | Ejecuta un comando, stdout en `#exec-log` |
| HTTP GET | `curl:https://url` | Petición HTTP, respuesta en `#exec-log` |
| Tema | `theme:nombre` / `theme:toggle` | Cambia tema activo |
| Suscripción | `subscribe:topic` | Suscribe stdin a un topic del event bus |

## Eventos del Sistema

El `EventBus` (Pub/Sub thread-safe) publica eventos en estos tópicos:

| Topic | Origen | Datos |
|-------|--------|-------|
| `stdin` | Pipe de entrada estándar | Línea de texto (string) |
| `exec-output` | HandlerExec | Output del comando |
| `theme-changed` | HandlerTheme | Nombre del nuevo tema |

## Atajos de Teclado

| Tecla | Acción |
|-------|--------|
| Tab | Siguiente elemento focusable |
| Shift+Tab | Elemento focusable anterior |
| ↑ ↓ ← → | Scroll |
| PgUp / PgDn | Scroll rápido (media página) |
| Home / End | Ir al inicio / final |
| Shift+D | Activar modo debug (colores aleatorios) |
| ESC | Cerrar modal |

## Tipos AST (Go)

```go
type Document struct {
    Pages []Page
}

type Page struct {
    Name     string
    Width, Height int
    Style    ComputedStyle
    Children []ElementNode
}

type ElementNode struct {
    Tag        string
    ID         string
    IncludeSrc string
    Attrs      map[string]string
    Events     map[string]string
    Style      ComputedStyle
    BoundBox   BoundBox
    Children   []ElementNode
    Text       string
}

type ComputedStyle struct {
    Width, Height         int
    MinWidth, MinHeight   int
    MaxWidth, MaxHeight   int
    Padding, Margin       int
    Grow                  float64
    Border                bool
    Direction             string   // "row" | "column"
    Gap                   int
    Align                 string   // "stretch" | "start" | "center" | "end"
    Justify               string   // "start" | "center" | "end" | "space-between"
    Overflow              string   // "visible" | "scroll"
    Color                 string
    BackgroundColor       string
    FontWeight            string   // "bold" | "normal"
}

type BoundBox struct {
    X, Y, W, H int
}
```

## Pipeline Interno

```
.hml → ParseHML() → AST raw → resolveIncludes()
→ ComputeStyles() (herencia + defaults + CSS vars)
→ Layout() (flexbox + min/max + padding/margin)
→ CellBuffer (Diff con frame anterior)
→ Screen.SetContent() → tcell.Show()
```

## Cobertura de Tests

| Paquete | Cobertura |
|---------|-----------|
| `internal/ast` | 100% |
| `internal/actions` | 100% |
| `internal/events` | 100% |
| `internal/theme` | 100% |
| `internal/parser` | 92.8% |
| `internal/layout` | 90.5% |
| `internal/render` | 79.1% |
