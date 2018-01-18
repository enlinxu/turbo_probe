package instance

import (
	"github.com/golang/glog"

	"turbo_probe/pkg/proto"
	"turbo_probe/pkg/probe"
)

type MockActionExecutor struct {
	name string
}

func NewMockActionExecutor(name string) *MockActionExecutor {
	return &MockActionExecutor{
		name: name,
	}
}

//3. Action
func (p *MockActionExecutor) ExecuteAction(actionExecutionDTO *proto.ActionExecutionDTO,
                                              accountValues []*proto.AccountValue,
                                              tracker probe.IActionProgressTracker) (*proto.ActionResult, error) {
	glog.Errorf("Action Execuator is not really implemented")

	state := proto.ActionResponseState_SUCCEEDED
	progress := int32(0)
	msg := "success"
	body := &proto.ActionResponse{
		ActionResponseState: &state,
		Progress: &progress,
		ResponseDescription: &msg,
	}

	result := &proto.ActionResult{
		Response: body,
	}
	return result, nil
}
