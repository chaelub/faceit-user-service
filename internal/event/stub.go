package event

import (
	"github.com/chaelub/faceit-user-service/internal/models"
	"log"
	"os"
	"sync"
)

type (
	eventHandler func([]byte) error
)

type StubEventService struct {
	log        *log.Logger
	handlersMu sync.RWMutex
	handlers   map[models.EventType][]eventHandler
}

func (es *StubEventService) Notify(evt models.Event) error {
	es.handlersMu.RLock()
	handlersList, got := es.handlers[evt.Type]
	if !got {
		es.handlersMu.RUnlock()
		return nil
	}
	handlers := make([]eventHandler, len(handlersList))
	copy(handlers, handlersList)
	es.handlersMu.RUnlock()

	for i := range handlers {
		err := handlers[i](evt.Message)
		if err != nil {
			return err
		}
	}

	return nil
}

func (es *StubEventService) RegisterHandler(evtType models.EventType, handler eventHandler) {
	es.handlersMu.Lock()
	es.handlers[evtType] = append(es.handlers[evtType], handler)
	es.handlersMu.Unlock()
}

func NewStubEventService() *StubEventService {
	return &StubEventService{
		handlers: make(map[models.EventType][]eventHandler),
		log:      log.New(os.Stderr, "[EventService] ", log.LstdFlags),
	}
}

func NewStubEventHandler(evtType models.EventType) eventHandler {
	logger := log.New(os.Stderr, "[StubEventHandler]", log.LstdFlags)
	return func(_ []byte) error {
		logger.Printf("Handler just received a new event with type %s", evtType.String())
		return nil
	}
}
