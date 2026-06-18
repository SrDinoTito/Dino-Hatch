// Package events tests para el bus pub/sub thread-safe.
package events

import (
	"sync"
	"testing"
	"time"
)

// recvTimeout recibe del canal con timeout. Usa 100ms por defecto.
func recvTimeout(t *testing.T, ch <-chan string, timeout time.Duration) (string, bool) {
	t.Helper()
	select {
	case val := <-ch:
		return val, true
	case <-time.After(timeout):
		return "", false
	}
}

// recvAll drena todas las lecturas disponibles de un canal y devuelve el count.
func recvAllCount(ch <-chan string) int {
	count := 0
	for {
		select {
		case <-ch:
			count++
		default:
			return count
		}
	}
}

// ---------------------------------------------------------------------------
// Test 1: Publish + Subscribe + ProcessAll básico
// ---------------------------------------------------------------------------

func TestPublishSubscribeBasic(t *testing.T) {
	bus := New()
	ch := bus.Subscribe("saludos")

	bus.Publish("saludos", "hola mundo")
	bus.ProcessAll()

	val, ok := recvTimeout(t, ch, 100*time.Millisecond)
	if !ok {
		t.Fatal("esperaba recibir mensaje pero no llegó nada")
	}
	if val != "hola mundo" {
		t.Errorf("esperaba 'hola mundo', recibí '%s'", val)
	}
}

// ---------------------------------------------------------------------------
// Test 2: Múltiples suscriptores al mismo tópico
// ---------------------------------------------------------------------------

func TestMultipleSubscribersSameTopic(t *testing.T) {
	bus := New()
	ch1 := bus.Subscribe("noticias")
	ch2 := bus.Subscribe("noticias")

	bus.Publish("noticias", "alerta")
	bus.ProcessAll()

	val1, ok1 := recvTimeout(t, ch1, 100*time.Millisecond)
	if !ok1 {
		t.Fatal("suscriptor 1 no recibió mensaje")
	}
	if val1 != "alerta" {
		t.Errorf("suscriptor 1 esperaba 'alerta', recibió '%s'", val1)
	}

	val2, ok2 := recvTimeout(t, ch2, 100*time.Millisecond)
	if !ok2 {
		t.Fatal("suscriptor 2 no recibió mensaje")
	}
	if val2 != "alerta" {
		t.Errorf("suscriptor 2 esperaba 'alerta', recibió '%s'", val2)
	}
}

// ---------------------------------------------------------------------------
// Test 3: Múltiples tópicos diferentes, cada suscriptor recibe solo el suyo
// ---------------------------------------------------------------------------

func TestMultipleTopicsIsolated(t *testing.T) {
	bus := New()
	chFoo := bus.Subscribe("foo")
	chBar := bus.Subscribe("bar")

	bus.Publish("foo", "mensaje-foo")
	bus.Publish("bar", "mensaje-bar")
	bus.ProcessAll()

	// Suscriptor foo recibe solo mensaje-foo
	valFoo, okFoo := recvTimeout(t, chFoo, 100*time.Millisecond)
	if !okFoo {
		t.Fatal("foo: no recibió mensaje")
	}
	if valFoo != "mensaje-foo" {
		t.Errorf("foo: esperaba 'mensaje-foo', recibió '%s'", valFoo)
	}

	// Suscriptor bar recibe solo mensaje-bar
	valBar, okBar := recvTimeout(t, chBar, 100*time.Millisecond)
	if !okBar {
		t.Fatal("bar: no recibió mensaje")
	}
	if valBar != "mensaje-bar" {
		t.Errorf("bar: esperaba 'mensaje-bar', recibió '%s'", valBar)
	}

	// Verificar que ningún suscriptor recibió mensajes extra
	select {
	case v := <-chFoo:
		t.Errorf("foo: mensaje inesperado '%s'", v)
	default:
	}
	select {
	case v := <-chBar:
		t.Errorf("bar: mensaje inesperado '%s'", v)
	default:
	}
}

// ---------------------------------------------------------------------------
// Test 4: Backpressure — buffer suscriptor lleno (64), mensaje se descarta.
// Publish no debe bloquearse.
// ---------------------------------------------------------------------------

func TestBackpressureDropsMessage(t *testing.T) {
	bus := New()
	_ = bus.Subscribe("test")

	// Publicar 65 mensajes; el buffer del suscriptor es 64.
	// Los primeros 64 entran, el 65 se descarta en ProcessAll.
	for i := 0; i < 65; i++ {
		bus.Publish("test", "data")
	}
	bus.ProcessAll()

	// Publicar otro mensaje — ProcessAll debe funcionar sin errores
	// incluso después de haber descartado el 65°.
	bus.Publish("test", "post-drop")
	bus.ProcessAll()
}

func TestBackpressureExactly64Received(t *testing.T) {
	bus := New()
	ch := bus.Subscribe("test")

	const n = 65
	for i := 0; i < n; i++ {
		bus.Publish("test", "data")
	}
	bus.ProcessAll()

	count := recvAllCount(ch)
	if count != 64 {
		t.Errorf("backpressure: esperaba 64 mensajes (buffer lleno), recibí %d", count)
	}
}

// ---------------------------------------------------------------------------
// Test 5: Thread-safety — publicaciones concurrentes sin data races.
// Se ejecuta con `go test -race` para detectar condiciones de carrera.
// ---------------------------------------------------------------------------

func TestConcurrentPublish(t *testing.T) {
	bus := New()
	ch := bus.Subscribe("test")

	var wg sync.WaitGroup
	goroutines := 10
	pubsPerGoroutine := 20
	totalPubs := goroutines * pubsPerGoroutine

	for g := 0; g < goroutines; g++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for i := 0; i < pubsPerGoroutine; i++ {
				bus.Publish("test", "data")
			}
		}(g)
	}
	wg.Wait()

	// Procesar todos los eventos publicados concurrentemente
	bus.ProcessAll()

	// Verificar que al menos algunos mensajes llegaron
	// (algunos pueden perderse por backpressure del buffer suscriptor)
	count := recvAllCount(ch)
	if count == 0 {
		t.Errorf("esperaba al menos 1 mensaje de %d publicaciones concurrentes, recibí 0", totalPubs)
	}
}

// TestConcurrentPublishMultipleTopics prueba thread-safety con múltiples tópicos.
func TestConcurrentPublishMultipleTopics(t *testing.T) {
	bus := New()
	chA := bus.Subscribe("topic-a")
	chB := bus.Subscribe("topic-b")

	var wg sync.WaitGroup
	pubsPerTopic := 20 // 2 topics × 2 goroutines × 20 = 80 < 256 (buffer publish)

	wg.Add(2)
	go func() {
		defer wg.Done()
		for i := 0; i < pubsPerTopic; i++ {
			bus.Publish("topic-a", "a")
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < pubsPerTopic; i++ {
			bus.Publish("topic-b", "b")
		}
	}()
	wg.Wait()

	bus.ProcessAll()

	// Verificar aislamiento: topic-a solo recibe "a", topic-b solo "b"
	for {
		select {
		case val := <-chA:
			if val != "a" {
				t.Errorf("topic-a recibió valor equivocado: '%s'", val)
			}
		default:
			goto checkB
		}
	}
checkB:
	for {
		select {
		case val := <-chB:
			if val != "b" {
				t.Errorf("topic-b recibió valor equivocado: '%s'", val)
			}
		default:
			return
		}
	}
}

// ---------------------------------------------------------------------------
// Test 6: Close — después de Close(), los canales suscriptores se cierran.
// ---------------------------------------------------------------------------

func TestCloseClosesSubscriberChannels(t *testing.T) {
	bus := New()
	ch := bus.Subscribe("test")

	bus.Close()

	// Leer del canal cerrado: el segundo valor (ok) debe ser false
	_, ok := <-ch
	if ok {
		t.Error("esperaba canal cerrado después de Close(), pero ok=true")
	}
}

func TestCloseAllSubscribersClosed(t *testing.T) {
	bus := New()
	ch1 := bus.Subscribe("a")
	ch2 := bus.Subscribe("a")
	ch3 := bus.Subscribe("b")

	bus.Close()

	_, ok1 := <-ch1
	_, ok2 := <-ch2
	_, ok3 := <-ch3

	if ok1 {
		t.Error("ch1 debería estar cerrado")
	}
	if ok2 {
		t.Error("ch2 debería estar cerrado")
	}
	if ok3 {
		t.Error("ch3 debería estar cerrado")
	}
}

// ---------------------------------------------------------------------------
// Edge cases adicionales
// ---------------------------------------------------------------------------

func TestProcessAllNoEvents(t *testing.T) {
	bus := New()
	// No hay eventos publicados — no debe paniquear ni colgarse
	bus.ProcessAll()
}

func TestProcessAllNoSubscribers(t *testing.T) {
	bus := New()
	bus.Publish("orphan", "nadie escucha")
	// No hay suscriptores — debe consumir el evento sin error
	bus.ProcessAll()
}

func TestProcessAllMultipleCalls(t *testing.T) {
	bus := New()
	ch := bus.Subscribe("test")

	bus.Publish("test", "uno")
	bus.ProcessAll()
	bus.Publish("test", "dos")
	bus.ProcessAll()

	val1, ok1 := recvTimeout(t, ch, 100*time.Millisecond)
	if !ok1 || val1 != "uno" {
		t.Errorf("esperaba 'uno', recibió '%s' (ok=%v)", val1, ok1)
	}
	val2, ok2 := recvTimeout(t, ch, 100*time.Millisecond)
	if !ok2 || val2 != "dos" {
		t.Errorf("esperaba 'dos', recibió '%s' (ok=%v)", val2, ok2)
	}
}

func TestPublishAfterCloseNoSubscribers(t *testing.T) {
	bus := New()
	bus.Subscribe("test")
	bus.Close()

	// Publicar después de Close no debe paniquear
	// (no hay subscriptores ya, pero publish channel sigue abierto)
	bus.Publish("test", "despues")
}
