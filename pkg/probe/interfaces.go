package probe

import (
	"turbo_probe/pkg/proto"
)

//1. registration information
type IProbeInfoProvider interface {
	GetProtocolVersion() string
	GetProbeInfo() ([]*proto.ProbeInfo, error)
}

//2. discovery
type IDiscoveryExecutor interface {
	Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error)
	//GetAccountValues() *TurboTargetInfo

	Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error)
	DiscoverIncremental(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error)
	DiscoverPerformance(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error)
}

// Interface to perform execution of an action request for an entity in the TurboProbe.
// It receives a ActionExecutionDTO that contains the action request parameters. The target account values contain the
// information for connecting to the target environment to which the entity belongs. ActionProgressTracker will be used
// by the client to send periodic action progress updates to the server.
type IActionExecutor interface {
	ExecuteAction(actionExecutionDTO *proto.ActionExecutionDTO,
		accountValues []*proto.AccountValue,
		progressTracker IActionProgressTracker) (*proto.ActionResult, error)
}

// Interface to send action progress to the server
type IActionProgressTracker interface {
	UpdateProgress(progress int32, description string)
}
