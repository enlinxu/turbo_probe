package mediation

import (
	"fmt"
	"turbo_probe/pkg/proto"

	"github.com/golang/glog"
	protobuf "github.com/golang/protobuf/proto"
)

/*
type iWorker interface {
	Handle(req *proto.MediationServerMessage) error
	Stop() error
}

type discoveryWorker struct {
	req *proto.MediationServerMessage
}*/

func (m *MediationClient) doDiscovery(req *proto.DiscoveryRequest) (*proto.DiscoveryResponse, error) {
	switch dtype := req.GetDiscoveryType(); dtype {
	case proto.DiscoveryType_FULL:
		glog.V(2).Infof("Begin to do full discovery: %v", req.GetProbeType())
		return m.probe.DiscoveryExecutor.Discover(req.GetAccountValue())
	case proto.DiscoveryType_INCREMENTAL:
		glog.V(2).Infof("Begin to do incremental discovery: %v", req.GetProbeType())
		return m.probe.DiscoveryExecutor.DiscoverIncremental(req.GetAccountValue())
	case proto.DiscoveryType_PERFORMANCE:
		glog.V(2).Infof("Begin to do performance discovery: %v", req.GetProbeType())
		return m.probe.DiscoveryExecutor.DiscoverPerformance(req.GetAccountValue())
	}

	err := fmt.Errorf("Unknown discovery type: %v", req.GetDiscoveryType())
	glog.Error(err.Error())
	return nil, err
}

func (m *MediationClient) handleDiscoveryInner(sreq *proto.MediationServerMessage) error {
	msgId := sreq.GetMessageID()
	glog.V(2).Infof("begin to handle discovery request(id=%d)", msgId)
	req := sreq.GetDiscoveryRequest()
	if req == nil {
		err := fmt.Errorf("Failed to get discovery request: %++v", sreq)
		glog.Error(err.Error())
		return err
	}

	//1. keep alive

	//2. do discovery with time out
	resp, err := m.doDiscovery(req)
	if err != nil {
		glog.Errorf("Discovery error: %v", err)
		return err
	}

	//3. send response
	m.sendDiscoveryResponse(msgId, resp)
	return nil
}

func (m *MediationClient) handleDiscovery(sreq *proto.MediationServerMessage) error {
	err := m.handleDiscoveryInner(sreq)
	if err == nil {
		return nil
	}

	level := proto.ErrorDTO_WARNING
	des := err.Error()

	errDTO := &proto.ErrorDTO{
		Severity:    &level,
		Description: &des,
	}

	resp := &proto.DiscoveryResponse{
		ErrorDTO: []*proto.ErrorDTO{errDTO},
	}
	m.sendDiscoveryResponse(sreq.GetMessageID(), resp)
	return nil
}

func (m *MediationClient) sendDiscoveryResponse(msgId int32, resp *proto.DiscoveryResponse) {
	msgbody := &proto.MediationClientMessage_DiscoveryResponse{
		DiscoveryResponse: resp,
	}
	msg := &proto.MediationClientMessage{
		MessageID:              &msgId,
		MediationClientMessage: msgbody,
	}

	dat, err := protobuf.Marshal(msg)
	if err != nil {
		glog.Errorf("Failed to marshal protobuf message: %+v", err)
	}
	m.wsconn.PushSend(dat, defaultSendTimeOut)
}

func (m *MediationClient) handleValidation(sreq *proto.MediationServerMessage) {
	err := m.handleValidationInner(sreq)
	if err == nil {
		return
	}

	level := proto.ErrorDTO_WARNING
	des := err.Error()

	errDTO := &proto.ErrorDTO{
		Severity:    &level,
		Description: &des,
	}

	resp := &proto.ValidationResponse{
		ErrorDTO: []*proto.ErrorDTO{errDTO},
	}

	m.sendValidationResponse(sreq.GetMessageID(), resp)
}

func (m *MediationClient) handleValidationInner(sreq *proto.MediationServerMessage) error {
	msgId := sreq.GetMessageID()
	glog.V(2).Infof("begin to handle validation request(id=%d)", msgId)
	req := sreq.GetValidationRequest()
	if req == nil {
		err := fmt.Errorf("Failed to get validation request: %+++v", sreq)
		glog.Error(err.Error())
		return err
	}

	resp, err := m.probe.DiscoveryExecutor.Validate(req.GetAccountValue())
	if err != nil {
		err = fmt.Errorf("req=%d, Failed to validate target: %v", msgId, err)
		glog.Errorf(err.Error())
		return err
	}

	m.sendValidationResponse(msgId, resp)

	return nil
}

func (m *MediationClient) sendValidationResponse(msgId int32, resp *proto.ValidationResponse) {
	msgbody := &proto.MediationClientMessage_ValidationResponse{
		ValidationResponse: resp,
	}
	msg := &proto.MediationClientMessage{
		MessageID:              &msgId,
		MediationClientMessage: msgbody,
	}

	dat, err := protobuf.Marshal(msg)
	if err != nil {
		glog.Errorf("Failed to marshal protobuf message: %+v", err)
	}
	m.wsconn.PushSend(dat, defaultSendTimeOut)
}
