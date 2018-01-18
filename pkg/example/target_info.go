package instance

import (
	"turbo_probe/pkg/proto"
)

// Abstraction for the TurboTarget object in the client
type TurboTargetInfo struct {
	// Category of the target, such as Hypervisor, Storage, etc
	targetCategory string
	// Type of the target, such as Kubernetes, vCenter, etc
	targetType string
	// The field that uniquely identifies the target.
	targetIdentifierField string
	// Account values, such as username, password, nameOrAddress, etc
	// NOTE, it should be consisted with the AccountDefEntry has been defined in registration client.
	accountValues []*proto.AccountValue
}
