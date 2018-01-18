package probe

import (
	"fmt"
)

type TurboProbe struct {
	ProbeInfoProvider IProbeInfoProvider
	DiscoveryExecutor IDiscoveryExecutor
	ActionExecutor    IActionExecutor
}

type TurboProbeBuilder struct {
	probe *TurboProbe

	hasProbeInfoProvider bool
	hasDiscoveryExecutor bool
	hasActionExecutor    bool
}

func NewTurboProbeBuilder() *TurboProbeBuilder {
	return &TurboProbeBuilder{
		probe:                &TurboProbe{},
		hasProbeInfoProvider: false,
		hasDiscoveryExecutor: false,
		hasActionExecutor:    false,
	}
}

func (b *TurboProbeBuilder) WithRegInfoProvider(reg IProbeInfoProvider) *TurboProbeBuilder {
	b.probe.ProbeInfoProvider = reg
	b.hasProbeInfoProvider = true
	return b
}

func (b *TurboProbeBuilder) WithDiscoveryExecutor(d IDiscoveryExecutor) *TurboProbeBuilder {
	b.probe.DiscoveryExecutor = d
	b.hasDiscoveryExecutor = true
	return b
}

func (b *TurboProbeBuilder) WithActionExecutor(a IActionExecutor) *TurboProbeBuilder {
	b.probe.ActionExecutor = a
	b.hasActionExecutor = true
	return b
}

func (b *TurboProbeBuilder) Create() (*TurboProbe, error) {
	if !b.hasProbeInfoProvider {
		return nil, fmt.Errorf("not have Registration Info Provider")
	}

	if !b.hasDiscoveryExecutor {
		return nil, fmt.Errorf("not has Discovery Executor")
	}

	if !b.hasActionExecutor {
		return nil, fmt.Errorf("not has Action Executor")
	}
	return b.probe, nil
}
