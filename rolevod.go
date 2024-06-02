package rolevod

import (
	"errors"
)

type Session interface {
	session()
}

type Data interface {
	data()
}

type Key string // локально уникально
type Ref string // глобально уникально
type Label string

type Variable[T any] struct {
	key  Key
	spec T
}

type Communicable[T any] struct {
	id       Ref
	spec     T
	messages []Message
}

type Copyable struct{}

type Message interface {
	message()
}

func (Communicable[T]) message() {}
func (Copyable) message()        {}

type Signature1 struct {
	consumables []Variable[Session]
	producible  Variable[Session]
}

type Signature2 struct {
	consumables map[string]Session
	producible  Session
}

type State struct {
	name string
}

type Process1 struct {
	id          Ref
	consumables map[string]*Communicable[Session]
	producible  Communicable[Session]
}

type Process2 struct {
	id          Ref
	consumables []Variable[Communicable[Session]]
	producible  Variable[Communicable[Session]]
}

type ProviderMessage[T any] struct {
	message  T
	behavior Session
}

func (ProviderMessage[T]) data()    {}
func (ProviderMessage[T]) session() {}

type ProviderChoice struct {
	branches map[Label]Session
}

func (ProviderChoice) session() {}

type ProviderClose struct{}

func (ProviderClose) data()    {}
func (ProviderClose) session() {}

type ConsumerMessage[T any] struct {
	message  T
	behavior Session
}

func (ConsumerMessage[T]) data()    {}
func (ConsumerMessage[T]) session() {}

type ConsumerChoice struct {
	branches map[Label]Session
}

func (ConsumerChoice) session() {}

type ConsumerClose struct{}

func (ConsumerClose) data()    {}
func (ConsumerClose) session() {}

type Rolevod struct {
	processes map[Ref]Process1
}

func (rv *Rolevod) Send1(sender Ref, receiver Key, output Key) error {
	process := rv.processes[sender]
	message, ok := process.consumables[string(output)]
	if !ok {
		return errors.New("No such receiver")
	}
	switch spec := process.producible.spec.(type) {
	case ProviderMessage[Session]:
		if spec.message != message.spec {
			return errors.New("Unexpected message spec")
		}
		messages := process.producible.messages
		if len(messages) > 0 {
			return errors.New("There are pending messages")
		}
		process.producible.messages = append(messages, message)
	case ProviderMessage[Data]:
	default:
		return errors.New("Unexpected spec type")
	}
	return nil
}

func (rv *Rolevod) Send3(from Ref, to Key, what Label) error {
	return nil
}
