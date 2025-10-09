package processors

import (
	"github.com/fgrzl/messaging"
	"github.com/fgrzl/streamkit"
	"github.com/google/uuid"
	"github.com/hydn-co/mesh-sdk/pkg/lifecycle"
)

type Processor interface {
	lifecycle.Service
	MessageBusProcessor
	StreamProcessor
}

func NewGlobalProcessorBase(
	bus messaging.MessageBus,
	stream streamkit.Client,
) *GlobalProcessorBase {
	return &GlobalProcessorBase{
		MessageBusProcessorBase: NewMessageBusProcessorBase(bus),
		StreamProcessorBase:     NewStreamProcessorBase(stream, uuid.Nil),
	}
}

type GlobalProcessorBase struct {
	*MessageBusProcessorBase
	*StreamProcessorBase
}

func NewScopedProcessorBase(
	bus messaging.MessageBus,
	stream streamkit.Client,
	storeID uuid.UUID,
) *ScopedProcessorBase {
	return &ScopedProcessorBase{
		MessageBusProcessorBase: NewMessageBusProcessorBase(bus),
		StreamProcessorBase:     NewStreamProcessorBase(stream, storeID),
	}
}

type ScopedProcessorBase struct {
	*MessageBusProcessorBase
	*StreamProcessorBase
}
