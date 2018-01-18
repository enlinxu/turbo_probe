package mediation

import (
	"github.com/golang/glog"
	protobuf "github.com/golang/protobuf/proto"
	"turbo_probe/pkg/proto"
)

type RequestType string

const (
	DISCOVERY_REQUEST  RequestType = "Discovery"
	VALIDATION_REQUEST RequestType = "Validation"
	INTERRUPT_REQUEST  RequestType = "Interrupt"
	ACTION_REQUEST     RequestType = "Action"
	PROPERTY_REQUEST   RequestType = "Property"
	UNKNOWN_REQUEST    RequestType = "Unknown"
)

func getRequestType(serverRequest *proto.MediationServerMessage) RequestType {
	if serverRequest.GetValidationRequest() != nil {
		return VALIDATION_REQUEST
	} else if serverRequest.GetDiscoveryRequest() != nil {
		return DISCOVERY_REQUEST
	} else if serverRequest.GetActionRequest() != nil {
		return ACTION_REQUEST
	} else if serverRequest.GetInterruptOperation() > 0 {
		return INTERRUPT_REQUEST
	} else if serverRequest.GetProperties() != nil {
		return PROPERTY_REQUEST
	} else {
		return UNKNOWN_REQUEST
	}
}

func (m *MediationClient) handleServerRequest(dat []byte) error {
	mrequest := &proto.MediationServerMessage{}

	if err := protobuf.Unmarshal(dat, mrequest); err != nil {
		glog.Errorf("Failed to unmarshal mediation server message: %v", dat)
		return err
	}

	mtype := getRequestType(mrequest)
	glog.V(2).Infof("Received mediation server msg: id=%d, type=%s", mrequest.GetMessageID(), mtype)

	//TODO: handle server request according to request time
	switch mtype {
	case DISCOVERY_REQUEST:
		m.handleDiscovery(mrequest)
	case VALIDATION_REQUEST:
		m.handleValidation(mrequest)
	case ACTION_REQUEST:
		m.handleAction(mrequest)
	case INTERRUPT_REQUEST:
		glog.Errorf("Received interruput request: %++v", mrequest.GetInterruptOperation())
	default:
		glog.Errorf("Unknow request: %v, %++v", mtype, mrequest)
	}
	return nil
}
