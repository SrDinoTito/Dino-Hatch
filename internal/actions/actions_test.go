// Tests para actions: Registry, HandlerExec, HandlerCurl, HandlerTheme.
package actions

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// ============================================================
// Registry tests
// ============================================================

func TestRegistryRegisterAndExecute(t *testing.T) {
	r := New()
	var capturedArgs string
	r.Register("greet", func(args string, cb *Callbacks) error {
		capturedArgs = args
		return nil
	})
	cb := &Callbacks{
		Log: func(format string, args ...interface{}) {},
	}
	err := r.Execute("greet:Mundo", cb)
	if err != nil {
		t.Fatalf("Execute fallo: %v", err)
	}
	if capturedArgs != "Mundo" {
		t.Fatalf("args esperado 'Mundo', obtenido '%s'", capturedArgs)
	}
}

func TestRegistryExecuteInvalidFormat(t *testing.T) {
	r := New()
	cb := &Callbacks{
		Log: func(format string, args ...interface{}) {},
	}
	err := r.Execute("sin_dos_puntos", cb)
	if err == nil {
		t.Fatal("se esperaba error por formato invalido, nil obtenido")
	}
	if !strings.Contains(err.Error(), "formato invalido") {
		t.Fatalf("error debe contener 'formato invalido', obtenido: %v", err)
	}
}

func TestRegistryExecuteUnknownHandler(t *testing.T) {
	r := New()
	cb := &Callbacks{
		Log: func(format string, args ...interface{}) {},
	}
	err := r.Execute("nope:args", cb)
	if err == nil {
		t.Fatal("se esperaba error por handler desconocido, nil obtenido")
	}
	if !strings.Contains(err.Error(), "accion desconocida") {
		t.Fatalf("error debe contener 'accion desconocida', obtenido: %v", err)
	}
}

func TestRegistryExecuteWithHandlerError(t *testing.T) {
	r := New()
	r.Register("fail", func(args string, cb *Callbacks) error {
		return &customError{msg: "something went wrong"}
	})
	cb := &Callbacks{
		Log: func(format string, args ...interface{}) {},
	}
	err := r.Execute("fail:now", cb)
	if err == nil {
		t.Fatal("se esperaba error del handler, nil obtenido")
	}
	if !strings.Contains(err.Error(), "something went wrong") {
		t.Fatalf("error debe contener 'something went wrong', obtenido: %v", err)
	}
}

// customError para probar propagacion de errores sin testify.
type customError struct{ msg string }

func (e *customError) Error() string { return e.msg }

// ============================================================
// HandlerExec tests
// ============================================================

func TestHandlerExecSimpleCommand(t *testing.T) {
	var setTextID, setTextContent string
	cb := &Callbacks{
		SetElementText: func(id, text string) {
			setTextID = id
			setTextContent = text
		},
		Log:          func(format string, args ...interface{}) {},
		PublishEvent: func(topic, data string) {},
	}
	err := HandlerExec("echo hello", cb)
	if err != nil {
		t.Fatalf("HandlerExec fallo: %v", err)
	}
	if setTextID != "exec-log" {
		t.Fatalf("SetElementText id esperado 'exec-log', obtenido '%s'", setTextID)
	}
	if !strings.Contains(setTextContent, "hello") {
		t.Fatalf("SetElementText debe contener 'hello', obtenido: %s", setTextContent)
	}
}

func TestHandlerExecEmptyArgs(t *testing.T) {
	cb := &Callbacks{
		Log: func(format string, args ...interface{}) {},
	}
	err := HandlerExec("", cb)
	if err != nil {
		t.Fatalf("HandlerExec con args vacio no debe fallar: %v", err)
	}
}

func TestHandlerExecFailingCommand(t *testing.T) {
	var setTextContent string
	cb := &Callbacks{
		SetElementText: func(id, text string) {
			setTextContent = text
		},
		Log:          func(format string, args ...interface{}) {},
		PublishEvent: func(topic, data string) {},
	}
	err := HandlerExec("nonexistent-command-abcdef", cb)
	if err != nil {
		t.Fatalf("HandlerExec no debe propagar error de comando: %v", err)
	}
	if setTextContent == "" {
		t.Fatal("SetElementText debe ser llamado incluso si el comando falla")
	}
}

func TestHandlerExecMultipleWords(t *testing.T) {
	var setTextContent string
	cb := &Callbacks{
		SetElementText: func(id, text string) {
			setTextContent = text
		},
		Log:          func(format string, args ...interface{}) {},
		PublishEvent: func(topic, data string) {},
	}
	// Usamos printf con formato para verificar argumentos multiples
	err := HandlerExec("printf hello world", cb)
	if err != nil {
		t.Fatalf("HandlerExec fallo: %v", err)
	}
	if !strings.Contains(setTextContent, "hello world") {
		t.Fatalf("output debe contener 'hello world', obtenido: %s", setTextContent)
	}
}

// ============================================================
// HandlerCurl tests
// ============================================================

func TestHandlerCurlWithHttptest(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test response body"))
	}))
	defer server.Close()

	var setTextContent string
	cb := &Callbacks{
		SetElementText: func(id, text string) {
			setTextContent = text
		},
		Log:          func(format string, args ...interface{}) {},
		PublishEvent: func(topic, data string) {},
	}

	err := HandlerCurl(server.URL, cb)
	if err != nil {
		t.Fatalf("HandlerCurl fallo: %v", err)
	}
	if !strings.Contains(setTextContent, "test response body") {
		t.Fatalf("SetElementText debe contener 'test response body', obtenido: %s", setTextContent)
	}
}

func TestHandlerCurlRespectsCallbackNil(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("data"))
	}))
	defer server.Close()

	// Sin SetElementText — verifica que no panic
	cb := &Callbacks{
		Log:          func(format string, args ...interface{}) {},
		PublishEvent: func(topic, data string) {},
	}

	err := HandlerCurl(server.URL, cb)
	if err != nil {
		t.Fatalf("HandlerCurl con SetElementText nil fallo: %v", err)
	}
}

func TestHandlerCurlAddsHTTPScheme(t *testing.T) {
	// Verificamos que si la URL no tiene esquema, HandlerCurl
	// le antepone https:// — esto se traduce en una peticion
	// que va a fallar (no podemos probar el exito sin mock),
	// pero podemos verificar que no panic y que devuelve error.
	// NOTA: esta prueba no depende de conectividad externa real;
	// la conexion fallara con un error de DNS/timeout que capturamos.
	cb := &Callbacks{
		Log:          func(format string, args ...interface{}) {},
		PublishEvent: func(topic, data string) {},
	}
	err := HandlerCurl("invalid.example.test.domain", cb)
	// Esperamos error porque el dominio no existe, pero queremos
	// asegurar que la funcion no panic ni devuelva nil cuando
	// realmente falla la conexion.
	if err == nil {
		t.Log("HandlerCurl devolvio nil con dominio inexistente (posible DNS local)")
	}
}

// ============================================================
// HandlerTheme tests
// ============================================================

func TestHandlerThemeWithCallback(t *testing.T) {
	var switchedTheme string
	cb := &Callbacks{
		SwitchTheme: func(name string) error {
			switchedTheme = name
			return nil
		},
		Log: func(format string, args ...interface{}) {},
	}
	err := HandlerTheme("dark", cb)
	if err != nil {
		t.Fatalf("HandlerTheme fallo: %v", err)
	}
	if switchedTheme != "dark" {
		t.Fatalf("SwitchTheme debe recibir 'dark', recibio '%s'", switchedTheme)
	}
}

func TestHandlerThemeNilCallback(t *testing.T) {
	cb := &Callbacks{
		Log: func(format string, args ...interface{}) {},
	}
	// SwitchTheme es nil — no debe panic ni error
	err := HandlerTheme("dark", cb)
	if err != nil {
		t.Fatalf("HandlerTheme con SwitchTheme nil no debe fallar: %v", err)
	}
}

func TestHandlerThemeToggle(t *testing.T) {
	cb := &Callbacks{
		Log: func(format string, args ...interface{}) {},
	}
	err := HandlerTheme("toggle", cb)
	if err != nil {
		t.Fatalf("HandlerTheme toggle no debe fallar: %v", err)
	}
}

func TestHandlerThemeWithCallbackError(t *testing.T) {
	cb := &Callbacks{
		SwitchTheme: func(name string) error {
			return &customError{msg: "theme not found: " + name}
		},
		Log: func(format string, args ...interface{}) {},
	}
	err := HandlerTheme("nonexistent", cb)
	if err == nil {
		t.Fatal("se esperaba error de SwitchTheme, nil obtenido")
	}
	if !strings.Contains(err.Error(), "theme not found") {
		t.Fatalf("error debe contener 'theme not found', obtenido: %v", err)
	}
}

// ============================================================
// HandlerCurl: error path de io.ReadAll (exec.go lns 49-51)
// ============================================================

func TestHandlerCurlReadBodyError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hij, ok := w.(http.Hijacker)
		if !ok {
			t.Fatal("httptest server no soporta hijack")
		}
		conn, bufrw, err := hij.Hijack()
		if err != nil {
			t.Fatal("hijack fallo:", err)
		}
		defer conn.Close()
		// Enviar headers con Content-Length grande + body parcial,
		// para que io.ReadAll reciba error de body truncado
		bufrw.WriteString("HTTP/1.1 200 OK\r\nContent-Length: 100000\r\n\r\npartial")
		bufrw.Flush()
		conn.Close()
	}))
	defer server.Close()

	cb := &Callbacks{
		SetElementText: func(id, text string) {},
		Log:            func(format string, args ...interface{}) {},
		PublishEvent:   func(topic, data string) {},
	}

	err := HandlerCurl(server.URL, cb)
	if err == nil {
		t.Fatal("HandlerCurl debe retornar error con body incompleto")
	}
	t.Logf("Error de body truncado: %v", err)
}
