package processors

import (
	"context"
	"fmt"

	"github.com/fgrzl/claims"
	"github.com/fgrzl/messaging"
	"github.com/google/uuid"
)

// RegisterMessageHandler is a generic helper for registering typed message handlers.
func RegisterMessageHandler[TMessage messaging.Message](
	p MessageBusProcessor,
	route messaging.Route,
	handler func(context.Context, TMessage) error,
) (messaging.Subscription, error) {
	wrapped := func(ctx context.Context, msg messaging.Message) error {
		cast, ok := msg.(TMessage)
		if !ok {
			return fmt.Errorf("unexpected message type: got %T", msg)
		}
		return handler(ctx, cast)
	}
	return p.RegisterMessageHandler(route, wrapped)
}

// RegisterRequestHandler is a generic helper to safely cast and register typed request handlers.
func RegisterRequestHandler[TRequest messaging.Request, TResponse messaging.Response](
	processor MessageBusProcessor,
	route messaging.Route,
	handler func(context.Context, TRequest) (TResponse, error),
) (messaging.Subscription, error) {
	wrapped := func(ctx context.Context, req messaging.Request) (messaging.Response, error) {
		castReq, ok := req.(TRequest)
		if !ok {
			return nil, fmt.Errorf("unexpected request type: got %T", req)
		}
		resp, err := handler(ctx, castReq)
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
	return processor.RegisterRequestHandler(route, wrapped)
}

// MessageBusProcessor defines a component capable of handling typed request-response interactions.
type MessageBusProcessor interface {
	// RegisterMessageHandler registers a handler for a specific route.
	RegisterMessageHandler(route messaging.Route, handler messaging.MessageHandler) (messaging.Subscription, error)

	// RegisterRequestHandler registers a handler for a specific route.
	RegisterRequestHandler(route messaging.Route, handler messaging.RequestHandler) (messaging.Subscription, error)
}

type MessageBusProcessorBase struct {
	SubscriptionTracker
	bus messaging.MessageBus
}

func NewMessageBusProcessorBase(bus messaging.MessageBus) *MessageBusProcessorBase {
	return &MessageBusProcessorBase{
		SubscriptionTracker: NewSubscriptionTracker(),
		bus:                 bus,
	}
}

func (p *MessageBusProcessorBase) RegisterMessageHandler(
	route messaging.Route,
	handler messaging.MessageHandler,
) (messaging.Subscription, error) {
	sub, err := p.bus.Subscribe(route, handler)
	if err != nil {
		return nil, err
	}
	p.Track(sub)
	return sub, nil
}

func (p *MessageBusProcessorBase) RegisterRequestHandler(
	route messaging.Route,
	handler messaging.RequestHandler,
) (messaging.Subscription, error) {
	sub, err := p.bus.SubscribeRequest(route, handler)
	if err != nil {
		return nil, err
	}
	p.Track(sub)
	return sub, nil
}

func (p *MessageBusProcessorBase) CurrentUser(ctx context.Context) (claims.Principal, bool) {
	return claims.UserFromContext(ctx)
}

func (p *MessageBusProcessorBase) CurrentUserID(ctx context.Context) (uuid.UUID, bool) {
	user, ok := p.CurrentUser(ctx)
	if !ok {
		return uuid.Nil, false
	}
	subject := user.Subject()
	id, err := uuid.Parse(subject)
	if err != nil {
		return uuid.Nil, false
	}
	return id, true
}
