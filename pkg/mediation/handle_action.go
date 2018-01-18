package mediation

import (
	"fmt"
	"github.com/golang/glog"
	protobuf "github.com/golang/protobuf/proto"
	"turbo_probe/pkg/proto"
)

func (m *MediationClient) handleAction(sreq *proto.MediationServerMessage) {
	msgId := sreq.GetMessageID()
	err := m.handleActionInner(sreq)
	if err == nil {
		glog.V(3).Infof("Action-%d successed.", msgId)
		return
	}

	glog.Errorf("Action-%d failed: %v", msgId, err)
	//send failure response
	state := proto.ActionResponseState_FAILED
	progress := int32(0)
	msg := err.Error()
	body := &proto.ActionResponse{
		ActionResponseState: &state,
		Progress:            &progress,
		ResponseDescription: &msg,
	}

	resp := &proto.ActionResult{
		Response: body,
	}

	m.sendActionResult(msgId, resp)
}

func (m *MediationClient) handleActionInner(sreq *proto.MediationServerMessage) error {
	//1. get action detail
	msgId := sreq.GetMessageID()
	glog.V(2).Infof("begin to handle validation request(id=%d)", msgId)
	req := sreq.GetActionRequest()
	if req == nil {
		err := fmt.Errorf("Failed to get action request: %+++v", sreq)
		glog.Error(err.Error())
		return err
	}

	//2. setup progress tracker
	tracker := NewProgressTracker()
	// track the progress of action
	stop := make(chan bool)
	defer close(stop)
	go func() {
		m.handleActionProgress(msgId, stop, tracker)
	}()

	//3. execute action TODO: add keep-alive
	resp, err := m.probe.ActionExecutor.ExecuteAction(req.GetActionExecutionDTO(), req.GetAccountValue(), tracker)
	if err != nil {
		glog.Errorf("Action-%d failed with %v", msgId, err)
		return err
	}
	m.sendActionResult(msgId, resp)

	return nil
}

func (m *MediationClient) handleActionProgress(msgId int32, stop chan bool, tracker *ProgressTracker) {
	msgCh := tracker.getMsgChan()

	for {
		select {
		case <-stop:
			glog.V(3).Infof("Stop progress handler for action-%d", msgId)
			return
		case info := <-msgCh:
			progress := info.Response.GetProgress()
			glog.V(4).Infof("Get action-%d progress info: %d, %v",
				msgId, progress, info.Response.GetResponseDescription())
			m.sendActionProgress(msgId, info)
			if progress >= 100 {
				glog.V(2).Infof("Action-%d progress = 100%", msgId)
				return
			}
		}
	}
}

func (m *MediationClient) sendActionResult(msgId int32, resp *proto.ActionResult) {
	msgbody := &proto.MediationClientMessage_ActionResponse{
		ActionResponse: resp,
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

func (m *MediationClient) sendActionProgress(msgId int32, resp *proto.ActionProgress) {
	msgbody := &proto.MediationClientMessage_ActionProgress{
		ActionProgress: resp,
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
