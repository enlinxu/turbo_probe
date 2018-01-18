package instance

import (
	"fmt"
	"github.com/golang/glog"
	"turbo_probe/pkg/proto"
)

type MockDiscoveryExecutor struct {
	name string
}

func NewMockDiscoveryExecutor(name string) *MockDiscoveryExecutor {
	return &MockDiscoveryExecutor{
		name: name,
	}
}

//2. for validation and 3 discovery functions
func (p *MockDiscoveryExecutor) Validate(accountValues []*proto.AccountValue) (*proto.ValidationResponse, error) {
	//TODO: implement it
	glog.Errorf("mocked validation, always return good.")
	result := &proto.ValidationResponse{}
	return result, nil
}

func (p *MockDiscoveryExecutor) Discover(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	//TODO: implement it
	glog.Errorf("mocked discovery, return empty")
	dtos := []*proto.EntityDTO{}
	result := &proto.DiscoveryResponse{
		EntityDTO: dtos,
	}
	return result, nil
}

func (p *MockDiscoveryExecutor) DiscoverIncremental(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	//TODO: implement it
	err := fmt.Errorf("DiscvoeryIncremental is not really implemented.")
	glog.Errorf(err.Error())
	return nil, err
}

func (p *MockDiscoveryExecutor) DiscoverPerformance(accountValues []*proto.AccountValue) (*proto.DiscoveryResponse, error) {
	//TODO: implement it
	err := fmt.Errorf("DiscvoeryPerformance is not really implemented.")
	glog.Error(err.Error())
	return nil, err
}
