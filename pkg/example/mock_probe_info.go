package instance

import (
	"fmt"
	"github.com/golang/glog"
	"turbo_probe/pkg/dtobuilder"
	"turbo_probe/pkg/proto"
)

const (
	TargetIdentifierField string = "targetIdentifier"
	Username              string = "username"
	Password              string = "password"
)

type MockProbeInfoProvider struct {
	protocolVersion              string
	probeType                    string
	probeCategory                string
	fullDiscoveryInterval        int32
	incrementalDiscoveryInterval int32
	performanceDiscoveryInterval int32
}

func NewMockProbeInfoProvider(ver, ptype, category string) *MockProbeInfoProvider {
	return &MockProbeInfoProvider{
		protocolVersion: ver,
		probeType:       ptype,
		probeCategory:   category,

		fullDiscoveryInterval:        int32(300),
		incrementalDiscoveryInterval: int32(-1),
		performanceDiscoveryInterval: int32(-1),
	}
}

//1. for registration
func (p *MockProbeInfoProvider) GetProtocolVersion() string {
	return p.protocolVersion
}

func (p *MockProbeInfoProvider) GetProbeInfo() ([]*proto.ProbeInfo, error) {
	//1. probe type & category
	result := []*proto.ProbeInfo{}

	info := &proto.ProbeInfo{
		ProbeType:     &p.probeType,
		ProbeCategory: &p.probeCategory,
	}

	//2. account def & ID field
	info.AccountDefinition = p.getAccountDefinition()
	info.TargetIdentifierField = []string{TargetIdentifierField}

	//3. supply chain
	supplychain, err := p.getSupplyChain()
	if err != nil {
		glog.Errorf("Failed to create ProbeInfo: %v", err)
		return result, nil
	}
	info.SupplyChainDefinitionSet = supplychain

	//4. Discovery intervals
	info.FullRediscoveryIntervalSeconds = &(p.fullDiscoveryInterval)
	if p.incrementalDiscoveryInterval > 0 {
		info.IncrementalRediscoveryIntervalSeconds = &(p.incrementalDiscoveryInterval)
	}
	if p.performanceDiscoveryInterval > 0 {
		info.PerformanceRediscoveryIntervalSeconds = &(p.performanceDiscoveryInterval)
	}

	result = append(result, info)
	return result, nil
}

func (p *MockProbeInfoProvider) getAccountDefinition() []*proto.AccountDefEntry {
	var acctDefProps []*proto.AccountDefEntry

	// target ID
	targetIDAcctDefEntry := dtobuilder.BuildAccountDefEntry(TargetIdentifierField, "Address",
		"IP of the Kubernetes master", ".*", false, false)
	acctDefProps = append(acctDefProps, targetIDAcctDefEntry)

	// username
	usernameAcctDefEntry := dtobuilder.BuildAccountDefEntry(Username, "Username",
		"Username of the Kubernetes master", ".*", false, false)
	acctDefProps = append(acctDefProps, usernameAcctDefEntry)

	// password
	passwordAcctDefEntry := dtobuilder.BuildAccountDefEntry(Password, "Password",
		"Password of the Kubernetes master", ".*", false, true)
	acctDefProps = append(acctDefProps, passwordAcctDefEntry)

	return acctDefProps
}

func (p *MockProbeInfoProvider) getSupplyChain() ([]*proto.TemplateDTO, error) {
	result := []*proto.TemplateDTO{}
	vCPU := proto.CommodityDTO_VCPU
	vMem := proto.CommodityDTO_VMEM

	vCpuTemplate := &proto.TemplateCommodity{CommodityType: &vCPU}
	vMemTemplate := &proto.TemplateCommodity{CommodityType: &vMem}

	//1. pm
	pmb := dtobuilder.NewTemplateDTOBuilder(proto.EntityDTO_PHYSICAL_MACHINE)
	pmb.Sells(vCpuTemplate).Sells(vMemTemplate)
	pm, err := pmb.Create()
	if err != nil {
		err = fmt.Errorf("Failed to create PM templateDTO: %v", err)
		glog.Errorf(err.Error())
		return result, err
	}
	result = append(result, pm)

	//2. vm
	vmb := dtobuilder.NewTemplateDTOBuilder(proto.EntityDTO_VIRTUAL_MACHINE)
	vmb.Sells(vCpuTemplate).Sells(vMemTemplate)
	vmb.Provider(proto.EntityDTO_PHYSICAL_MACHINE, proto.Provider_HOSTING)
	vmb.Buys(vCpuTemplate).Buys(vMemTemplate)
	vm, err := vmb.Create()
	if err != nil {
		err = fmt.Errorf("Failed to create vm templateDTO: %v", err)
		glog.Errorf(err.Error())
		return result, err
	}

	result = append(result, vm)
	return result, nil
}
