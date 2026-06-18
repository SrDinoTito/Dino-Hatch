# HATCH-CORE — Requirements

## Version
0.1.0-alpha

## Objetivo
Construir el núcleo del framework Hatch: un pipeline completo que parsea archivos .hml
(HML + HSS) en un AST, resuelve includes de componentes, aplica estilos computados con
CSS variables, calcula layout flexbox con padding/margin/min/max clamping, soporta
overflow scroll por contenedor, eventos declarativos y data binding, y renderiza la
interfaz en terminal usando tcell v3 con cell buffer persistente y diff optimizado.

## Contexto
No existe un framework TUI declarativo en Go. Bubble Tea es el estándar de facto
(Elm Architecture, MVU) pero requiere layout manual. Hatch propone un enfoque
declarativo donde el usuario describe la UI en archivos .hml (TermML) con estilos HSS,
y Hatch se encarga del layout y renderizado automático con componentes reutilizables,
scroll, eventos y theming.

## Alcance
MVP (v0.1.0-alpha): Parseo HML/HSS con componentes (<include>), layout flexbox con
padding/margin/min/max, overflow scroll por contenedor, cell buffer persistente con
dirty flag, scrollbar visual, eventos declarativos (onclick/onchange/onfocus/onblur),
data binding (bind), CSS variables (:root + var()), tab focus navigation, y demo con
4 páginas y 7 componentes. No incluye: hot-reload, animaciones, layouts grid,
soporte multi-monitor, accesibilidad.

---

## Requerimientos Funcionales

### REQ-001 — Parseo de HML
El sistema debe parsear un archivo .hml con sintaxis XML-like en un AST. El formato
incluye etiquetas como `<page>`, `<box>`, `<text>`, `<style>`, `<input>`, `<textarea>`,
`<button>`, `<include>`.
- **Criterio de aceptación**: Dado un .hml válido, se produce un AST con
  Document → Page → ElementNode. Atributos onclick, onchange, onfocus, onblur se
  extraen de Attrs a Events.
- **Edge cases**: Etiquetas sin cerrar → error descriptivo. Atributos malformados →
  error. Archivo vacío → error. Tag `<include>` con src → almacena IncludeSrc.
- **Prioridad**: Alta

### REQ-002 — Parseo de HSS
El sistema debe parsear bloques `<style>` dentro del .hml que contengan reglas CSS-like
(selectores, propiedades: valor). El selector `:root` se ignora en StyleRules pero se
procesa con ParseCSSVars para extraer variables CSS.
- **Criterio de aceptación**: Dado `<style> box { width: 100; height: 50; } </style>`,
  se produce un slice de StyleRule. Propiedades desconocidas → ignorar con warning.
- **Edge cases**: Propiedades desconocidas → ignorar con warning. Selectores no usados
  → ignorar soporte. `:root` no genera StyleRule.
- **Prioridad**: Alta

### REQ-003 — AST con estilos computados
El AST debe contener nodos con estilos computados (resolución de herencia y defaults).
Incluye propiedades: width, height, grow, gap, align, justify, direction, color, bg,
border, text-align, valign, padding, margin, min-width, min-height, max-width,
max-height, overflow.
- **Criterio de aceptación**: Cada ElementNode tiene un ComputedStyle completo con
  valores resueltos. Herencia: hijo hereda color, bg, direction, align, justify,
  text-align, valign, overflow del padre si no los define explícitamente.
- **Prioridad**: Alta

### REQ-004 — Layout flexbox
El sistema debe calcular posiciones (x, y, width, height) para cada elemento usando
un modelo flexbox simplificado con: grow, gap, align-items, justify-content, direction,
padding, margin, min-width/min-height, max-width/max-height.
- **Criterio de aceptación**: Dado un AST con estilos, se produce un árbol con BoundBox
  (x, y, w, h) para cada nodo. Padding/margin se aplican. Min/max clamping se ejecutan.
- **Edge cases**: Overflow (hijos suman más que el padre) → grow no recibe espacio.
  Contenedor sin hijos → ocupa 0. Post-expansion de contenedores sin width/height
  explícitos.
- **Prioridad**: Alta

### REQ-005 — Renderizado a terminal
El sistema debe dibujar el layout calculado en una terminal real usando tcell v3.
Soporta tags: box, text, button, input, textarea. Border, color, bg, text-align,
valign.
- **Criterio de aceptación**: Tras llamar a render, se ve la interfaz en terminal.
  Texto, colores, bordes se muestran correctamente.
- **Prioridad**: Alta

### REQ-006 — Cell buffer con diff persistente
El sistema debe usar un cell buffer intermedio persistente (PrevCB/CurrCB) y solo
enviar a terminal las celdas que cambiaron entre frames (diff). Incluye dirty flag
que permite saltar layout+render cuando no hay cambios.
- **Criterio de aceptación**: En el segundo frame sin cambios, no se ejecuta layout
  ni render. En frame con cambios, solo se actualizan las celdas modificadas.
- **Prioridad**: Media

### REQ-007 — Línea de comandos
El sistema debe tener un CLI que acepte `hatch run <archivo.hml>`.
- **Criterio de aceptación**: `hatch run canva/demo.hml` renderiza el archivo.
- **Prioridad**: Media

### REQ-013 — Component system (include)
El sistema debe soportar la inclusión de componentes reutilizables mediante la
etiqueta `<include src="ruta">`. Los includes se resuelven después del parseo HML
y antes de ComputeStyles. Soportan atributos merge (los del include ganan sobre los
del componente) e includes anidados.
- **Criterio de aceptación**: Dado `<include src="components/header.hml">`, el
  contenido del componente se resuelve y reemplaza al nodo include en el AST.
  Los estilos del componente (bloques `<style>`) se mergean en las reglas globales.
- **Edge cases**: Include anidado (componente incluye otro componente). Include con
  id → el id se asigna al root del componente. Archivo no encontrado → warning, no crash.
- **Prioridad**: Alta

### REQ-014 — Overflow scroll en contenedores
El sistema debe soportar scroll por contenedor mediante la propiedad CSS `overflow`
con valores "visible", "hidden", "scroll". Los contenedores con overflow:scroll
tienen un ScrollY interno independiente del scroll global de página.
- **Criterio de aceptación**: Hijos de un contenedor overflow:scroll se renderizan
  con clipping (isChildVisible/clipHeight). La rueda del mouse sobre el contenedor
  scrollea internamente. Se dibuja una scrollbar visual en el margen derecho del
  contenedor.
- **Edge cases**: Contenedor scrollable anidado dentro de otro contenedor scrollable.
  Hijos parcialmente visibles (top/bottom clip). Scrollbar proporcional.
- **Prioridad**: Alta

### REQ-015 — Eventos declarativos
El sistema debe soportar eventos declarativos en elementos mediante atributos
onclick, onchange, onfocus, onblur. El formato de acción es `tipo:argumento`.
- **Acciones soportadas**:
  - `page:NOMBRE` — Navega a la página especificada (ej: onclick="page:proyectos")
  - `modal:open` / `modal:close` / `modal:toggle` — Abre/cierra modal
  - `action:random_colors` / `action:quit` — Acciones globales
- **Criterio de aceptación**: Click en botón con onclick="page:proyectos" carga la
  página proyectos. onchange se dispara en input/textarea.
- **Prioridad**: Alta

### REQ-016 — Data binding
El sistema debe soportar data binding mediante el atributo `bind` en elementos
input/textarea. bind="target-id" actualiza el texto del elemento con ese id cada
vez que el input cambia.
- **Criterio de aceptación**: Escribir en un input con bind="output-text" actualiza
  el texto del elemento con id "output-text" en tiempo real.
- **Edge cases**: Target no encontrado → ignorar silenciosamente. Múltiples inputs
  bindeando al mismo target → el último cambio gana.
- **Prioridad**: Alta

### REQ-017 — CSS variables / themes
El sistema debe soportar CSS variables definidas en el selector `:root` y referenciadas
con `var(--nombre)` en propiedades de estilo.
- **Criterio de aceptación**: Dado `:root { --bg: #333; }` y `box { bg: var(--bg); }`,
  los boxes heredan bg=#333. Las variables se resuelven en applyProps.
- **Edge cases**: Variable no definida → se deja var(--nombre) sin resolver (ignorado).
  Variables globales se propagan vía ComputeStyles.
- **Prioridad**: Media

### REQ-018 — Tab focus navigation
El sistema debe soportar navegación por Tab entre elementos focusables (input,
textarea, button). El orden se construye por DFS. Shift+Tab navega en reversa.
- **Criterio de aceptación**: Tab mueve el foco al siguiente elemento focusable.
  Shift+Tab al anterior. El cursor aparece en input/textarea al recibir foco.
- **Edge cases**: Sin elementos focusables → Tab no hace nada. FocusIndex cíclico
  (vuelve al inicio al llegar al final).
- **Prioridad**: Media

---

## Requerimientos No Funcionales

### REQ-008 — Dependencias mínimas
Solo dependencias externas permitidas: tcell v3 + librerías de stdlib (encoding/xml,
encoding/json).
- **Prioridad**: Alta

### REQ-009 — Comentarios en español
Todo comentario en el código debe estar en español, ser mínimo y solo explicar el
"por qué", no el "qué".
- **Prioridad**: Media

### REQ-010 — Límite de lineas por archivo
Ningún archivo .go debe superar las 150 líneas. Si un archivo excede, se divide en
una subcarpeta con múltiples archivos. Excepción: layout.go (213L) por ser el core
algorítmico del layout engine.
- **Prioridad**: Media

### REQ-011 — Tests obligatorios
Todo paquete debe tener tests. Cobertura mínima: 70% en parser, 70% en layout,
70% en render, 70% en ast.
- **Prioridad**: Alta

### REQ-012 — Commits versionados
Los mensajes de commit deben seguir el formato `[0.1.X] descripción`. X se
incrementa en 1 por cada commit.
- **Prioridad**: Media

---

## Dependencias Externas
- **tcell v3** (`github.com/gdamore/tcell/v2` v2.13.10) — terminal cell buffer,
  eventos, resize
- **encoding/xml** (stdlib) — parseo de HML
- **encoding/json** (stdlib) — parseo de handler.json

## Riesgos
1. **Parser XML frágil**: HML no es XML puro — requiere preprocesamiento de `<style>`.
   Mitigación: preprocessStyleBlocks separa contenido CSS antes del parseo XML.
2. **Rendimiento de diff**: CellBuffer con diff en full-screen puede ser lento.
   Mitigación: dirty flag + throttle de mouse (33ms).
3. **Layout overflow scroll**: Hijos que escapan del viewport del contenedor deben
   clippearse correctamente. Mitigación: isChildVisible + clipHeight.
4. **Include anidados**: La resolución recursiva de includes puede causar ciclos.
   Mitigación: límite de profundidad implícito por recursión controlada.

## Módulos Afectados
- `internal/ast/` — tipos del AST (Document, Page, ElementNode, ComputedStyle, etc.)
- `internal/parser/` — parser HML + HSS, compute, props, inherit
- `internal/layout/` — engine flexbox + content height/width
- `internal/render/` — cell buffer + tcell driver
- `internal/handler/` — engine de bindings de teclado
- `cmd/hatch/` — CLI entry point (22 archivos: pipeline, eventloop, render, etc.)
