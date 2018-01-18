package instance

import (
	"fmt"
	"testing"
)

func TestMockProbe_GetProbeInfo(t *testing.T) {
	p := NewMockProbeInfoProvider("6.0.1", "kubernetes", "cloudnative")
	chain, err := p.getSupplyChain()
	if err != nil {
		t.Errorf("Failed to create supplychain: %v", err)
	}

	fmt.Printf("protocol version: %v\n", p.protocolVersion)

	for i := range chain {
		fmt.Printf("supply chain node: %++v\n", chain[i])
	}
}
