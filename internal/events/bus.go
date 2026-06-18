// Package events implementa un bus publicador/suscriptor thread-safe
// para comunicacion entre goroutines background y el event loop principal.
package events

import "sync"

// Event es un evento del bus con un topico y datos opcionales.
type Event struct {
	Topic string
	Data  string
}

// Bus es un publicador/suscriptor de eventos thread-safe.
// Los eventos se encolan y se consumen en el hilo principal (event loop).
type Bus struct {
	mu      sync.Mutex
	subs    map[string][]chan string
	publish chan Event // canal para recibir eventos publicados
	quit    chan struct{}
}

// New crea un bus con buffer de publicacion.
func New() *Bus {
	return &Bus{
		subs:    make(map[string][]chan string),
		publish: make(chan Event, 256),
		quit:    make(chan struct{}),
	}
}

// Subscribe crea un canal para recibir eventos de un topico.
// El canal tiene buffer 64 para no bloquear al publicador.
func (b *Bus) Subscribe(topic string) <-chan string {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan string, 64)
	b.subs[topic] = append(b.subs[topic], ch)
	return ch
}

// Publish encola un evento para ser distribuido.
// Es thread-safe, llamable desde cualquier goroutine.
func (b *Bus) Publish(topic, data string) {
	b.publish <- Event{Topic: topic, Data: data}
}

// ProcessAll distribuye todos los eventos encolados a sus suscriptores.
// Debe llamarse desde el event loop (hilo principal) periodica o inmediatamente.
func (b *Bus) ProcessAll() {
	b.mu.Lock()
	defer b.mu.Unlock()
	for {
		select {
		case ev := <-b.publish:
			for _, ch := range b.subs[ev.Topic] {
				select {
				case ch <- ev.Data:
				default:
					// Canal lleno, descartar evento (backpressure)
				}
			}
		default:
			return // no hay mas eventos
		}
	}
}

// Close cierra el bus y todos los canales suscriptores.
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	close(b.quit)
	for _, chans := range b.subs {
		for _, ch := range chans {
			close(ch)
		}
	}
	b.subs = make(map[string][]chan string)
}
