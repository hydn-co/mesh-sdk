package processors

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/fgrzl/collections"
	"github.com/fgrzl/enumerators"
	"github.com/fgrzl/json/polymorphic"
	"github.com/fgrzl/lexkey"
	"github.com/fgrzl/streamkit"
	"github.com/fgrzl/streamkit/pkg/api"
	"github.com/fgrzl/tickle"
	"github.com/google/uuid"
)

// PolymorphicStreamHandler is the signature for registered stream handlers.
type PolymorphicStreamHandler = func(context.Context, api.Consumable) error

// StreamProcessor defines the interface required to register handlers.
type StreamProcessor interface {
	RegisterStreamHandler(discriminator string, handler PolymorphicStreamHandler)
	RegisterSpaces(spaces ...string)
	RegisterOffsetLoader(func(ctx context.Context) (*ConsumerOffset, error))
	RegisterFlushOffset(func(ctx context.Context, off *ConsumerOffset) error)
}

// RegisterStreamHandler binds a strongly typed handler for a specific stream event.
func RegisterStreamHandler[T api.Consumable](p StreamProcessor, handler func(context.Context, T) error) error {
	var zero T
	discriminator := zero.GetDiscriminator()
	spaces := zero.GetSpaces()

	wrapper := func(ctx context.Context, consumable api.Consumable) error {
		content, ok := consumable.(T)
		if !ok {
			return errors.New("failed to cast content")
		}
		return handler(ctx, content)
	}

	p.RegisterSpaces(spaces...)
	p.RegisterStreamHandler(discriminator, wrapper)

	return nil
}

func RegisterOffsetLoader(p StreamProcessor, loader func(ctx context.Context) (*ConsumerOffset, error)) {
	p.RegisterOffsetLoader(loader)
}

func RegisterFlushOffset(p StreamProcessor, flush func(ctx context.Context, off *ConsumerOffset) error) {
	p.RegisterFlushOffset(flush)
}

// ConsumerOptions configure the stream consumer at startup.
type ConsumerOptions func(*StreamProcessorBase)

func WithOffset(offset *ConsumerOffset) ConsumerOptions {
	return func(p *StreamProcessorBase) {
		p.offset = offset
	}
}

func WithFlushHook(hook func(context.Context, *ConsumerOffset) error) ConsumerOptions {
	return func(p *StreamProcessorBase) {
		p.flushOffset = hook
	}
}

func WithBatchSize(n int) ConsumerOptions {
	return func(p *StreamProcessorBase) {
		p.batchSize = n
	}
}

// StreamProcessorBase is a reusable stream processor implementation.
type StreamProcessorBase struct {
	stream         streamkit.Client
	tickler        *tickle.Tickler
	storeID        uuid.UUID
	spaces         *collections.HashSet[string]
	subs           []api.Subscription
	streamHandlers map[string]PolymorphicStreamHandler
	loadOffset     func(ctx context.Context) (*ConsumerOffset, error)
	flushOffset    func(context.Context, *ConsumerOffset) error
	offset         *ConsumerOffset
	batchSize      int
	// lifecycle controls for the consumer goroutine
	mu        sync.Mutex
	runCancel context.CancelFunc
	runDone   chan struct{}
	running   bool
}

// NewStreamProcessorBase creates a new base processor with sensible defaults.
func NewStreamProcessorBase(stream streamkit.Client, storeID uuid.UUID) *StreamProcessorBase {
	return &StreamProcessorBase{
		stream:         stream,
		tickler:        tickle.NewTickler(),
		storeID:        storeID,
		spaces:         collections.NewHashSet[string](),
		streamHandlers: make(map[string]PolymorphicStreamHandler),
	}
}

// RegisterStreamHandler registers a handler for a specific event discriminator.
func (p *StreamProcessorBase) RegisterStreamHandler(discriminator string, handler PolymorphicStreamHandler) {

	if _, exists := p.streamHandlers[discriminator]; exists {
		panic(fmt.Sprintf("handler already registered for discriminator %q", discriminator))
	}

	p.streamHandlers[discriminator] = handler
}

func (p *StreamProcessorBase) RegisterSpaces(spaces ...string) {
	for _, space := range spaces {
		p.spaces.Add(space)
	}
}

func (p *StreamProcessorBase) RegisterOffsetLoader(load func(context.Context) (*ConsumerOffset, error)) {
	p.loadOffset = load
}

func (p *StreamProcessorBase) RegisterFlushOffset(flush func(context.Context, *ConsumerOffset) error) {
	p.flushOffset = flush
}

// StartConsumer begins consuming events from registered stream spaces.
//
// StartConsumer is intentionally non-blocking: it applies any provided options
// and launches the consumer loop in a new goroutine. The blocking work is
// performed in runConsumer(ctx).
func (p *StreamProcessorBase) StartConsumer(ctx context.Context, opts ...ConsumerOptions) error {
	if p.spaces.Size() == 0 {
		return nil
	}

	for _, opt := range opts {
		opt(p)
	}

	p.mu.Lock()
	if p.running {
		p.mu.Unlock()
		return nil
	}
	// create cancellable child context so Stop can cancel the consumer
	childCtx, cancel := context.WithCancel(ctx)
	p.runCancel = cancel
	p.running = true
	p.runDone = make(chan struct{})
	p.mu.Unlock()

	go func() {
		defer func() {
			// signal completion
			p.mu.Lock()
			if p.runDone != nil {
				close(p.runDone)
			}
			p.running = false
			p.runCancel = nil
			p.runDone = nil
			p.mu.Unlock()
		}()
		if err := p.runConsumer(childCtx); err != nil {
			// log to stdout for now; projects can swap this for structured logging
			fmt.Printf("stream processor run error: %v\n", err)
		}
	}()
	return nil
}

// Stop cancels a running consumer (if any) and waits for it to finish. It
// respects the provided context for timeout/cancellation while waiting.
func (p *StreamProcessorBase) Stop(ctx context.Context) error {
	p.mu.Lock()
	if !p.running || p.runCancel == nil {
		p.mu.Unlock()
		return nil
	}
	// grab cancel and waitgroup, then unlock before canceling to avoid deadlock
	cancel := p.runCancel
	p.mu.Unlock()

	// request shutdown
	cancel()

	// wait for run goroutine to exit or ctx cancellation
	p.mu.Lock()
	done := p.runDone
	p.mu.Unlock()

	if done == nil {
		return nil
	}

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// runConsumer contains the blocking consumer loop previously in StartConsumer.
func (p *StreamProcessorBase) runConsumer(ctx context.Context) error {
	// If a loader was registered, try to load previously persisted offsets.
	if p.loadOffset != nil {
		off, err := p.loadOffset(ctx)
		if err != nil {
			return fmt.Errorf("failed to load consumer offset: %w", err)
		}
		if off != nil {
			p.offset = off
		}
	}

	if p.offset == nil {
		p.offset = &ConsumerOffset{Offsets: make(map[string]lexkey.LexKey)}
	} else if p.offset.Offsets == nil {
		p.offset.Offsets = make(map[string]lexkey.LexKey)
	}

	// Default to flushing every event if batchSize not specified
	if p.batchSize == 0 {
		p.batchSize = 1
	}

	spaces := p.spaces.ToSlice()
	for _, space := range spaces {
		sub, err := p.stream.SubscribeToSpace(ctx, p.storeID, space, p.handleSegmentStatus)
		if err != nil {
			return fmt.Errorf("failed to subscribe to space %q: %w", space, err)
		}
		p.subs = append(p.subs, sub)

		if _, ok := p.offset.Offsets[space]; !ok {
			p.offset.Offsets[space] = lexkey.Empty
		}
	}

	sub := p.tickler.Subscribe(ctx, spaces...)
	defer sub.Dispose()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		default:
			var counter int
			args := &streamkit.Consume{Offsets: p.offset.Offsets}
			enumerator := p.stream.Consume(ctx, p.storeID, args)

			err := enumerators.ForEach(enumerator, func(entry *streamkit.Entry) error {
				if err := p.handleEntry(ctx, entry); err != nil {
					return err
				}

				// update offset for the space so next Consume resumes after this entry
				if entry != nil && entry.Space != "" {
					p.offset.Offsets[entry.Space] = entry.GetSpaceOffset()
				}

				counter++
				if p.flushOffset != nil && counter >= p.batchSize {
					if err := p.flushOffset(ctx, p.offset); err != nil {
						return err
					}
					counter = 0
				}
				return nil
			})
			if err != nil {
				fmt.Printf("error handling entries: %v\n", err)
			}

			if p.flushOffset != nil && counter > 0 {
				if err := p.flushOffset(ctx, p.offset); err != nil {
					return err
				}
			}

			sub.WaitTimeout(5 * time.Minute)
		}
	}
}

func (p *StreamProcessorBase) handleEntry(ctx context.Context, entry *streamkit.Entry) error {
	envelope, err := polymorphic.UnmarshalPolymorphicJSON(entry.Payload)
	if err != nil {
		return err
	}

	handler, ok := p.streamHandlers[envelope.Discriminator]
	if !ok {
		return fmt.Errorf("no handler registered for discriminator %q", envelope.Discriminator)
	}

	content, ok := envelope.Content.(api.Consumable)
	if !ok {
		return fmt.Errorf("invalid content type: %T", envelope.Content)
	}

	return handler(ctx, content)
}

func (p *StreamProcessorBase) handleSegmentStatus(status *streamkit.SegmentStatus) {
	p.tickler.Tickle(status.Space)
}

// ConsumerOffset tracks the per-space stream offsets.
type ConsumerOffset struct {
	Offsets map[string]lexkey.LexKey `json:"offsets"`
}
