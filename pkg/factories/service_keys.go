package factories

// ServiceKeys contains constants for dependency injection service keys
// used throughout the web application for consistent service registration
// and retrieval. These keys are used with mux.WithService() and mux.GetService().
//
// When adding new service keys:
// 1. Use descriptive names ending with "Key"
// 2. Ensure the value matches the interface or factory type name
// 3. Add appropriate documentation
const (
	// MessagingBusFactoryKey is the service key for registering and retrieving
	// the messaging bus factory dependency used for message-based communication.
	MessagingBusFactoryKey = "MessagingBusFactory"

	// StreamkitClientFactoryKey is the service key for registering and retrieving
	// the messaging streamkit client factory dependency used for message-based communication.
	StreamkitClientFactoryKey = "StreamkitClientFactoryKey"
)
