package testkit

import (
	"context"
	"time"

	"github.com/fgrzl/messaging"
	"github.com/stretchr/testify/mock"
)

// NewMockMessageBus instantiates a new mock message bus that supports stubbing and assertions.
func NewMockMessageBus() *MockMessageBus {
	return &MockMessageBus{}
}

// ConfigureBusFactory lets you specify the bus the factory should return
func ConfigureBusFactory(bus messaging.MessageBus) messaging.MessageBusFactory {
	return &MockBusFactory{
		bus: bus,
	}
}

type MockBusFactory struct {
	bus messaging.MessageBus
}

func (f *MockBusFactory) Get(ctx context.Context) (messaging.MessageBus, error) {
	return f.bus, nil
}

type MockMessageBus struct {
	mock.Mock
}

func (m *MockMessageBus) Notify(msg messaging.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

func (m *MockMessageBus) NotifyWithContext(ctx context.Context, msg messaging.Message) error {
	args := m.Called(ctx, msg)
	return args.Error(0)
}

func (m *MockMessageBus) Request(msg messaging.Request, timeout time.Duration) (messaging.Response, error) {
	args := m.Called(msg, timeout)
	return args.Get(0).(messaging.Response), args.Error(1)
}

func (m *MockMessageBus) RequestWithContext(ctx context.Context, msg messaging.Request, timeout time.Duration) (messaging.Response, error) {
	args := m.Called(ctx, msg, timeout)

	response := args.Get(0)
	err := args.Error(1)

	return response.(messaging.Response), err
}

func (m *MockMessageBus) Subscribe(route messaging.Route, handler messaging.MessageHandler) (messaging.Subscription, error) {
	args := m.Called(route, handler)
	return args.Get(0).(messaging.Subscription), args.Error(1)
}

func (m *MockMessageBus) SubscribeRequest(route messaging.Route, handler messaging.RequestHandler) (messaging.Subscription, error) {
	args := m.Called(route, handler)
	return args.Get(0).(messaging.Subscription), args.Error(1)
}

func (m *MockMessageBus) Close() error {
	args := m.Called()
	return args.Error(0)
}
