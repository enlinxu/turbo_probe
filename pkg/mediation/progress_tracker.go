package mediation

import (
	"turbo_probe/pkg/proto"
)

type ProgressTracker struct {
	msgCh chan *proto.ActionProgress
}

func NewProgressTracker() *ProgressTracker {
	return &ProgressTracker{
		msgCh: make(chan *proto.ActionProgress),
	}
}

func newProgress(progress int32, des string) *proto.ActionProgress {
	state := proto.ActionResponseState_IN_PROGRESS
	resp := &proto.ActionResponse{
		ActionResponseState: &state,
		Progress:            &progress,
		ResponseDescription: &des,
	}

	return &proto.ActionProgress{
		Response: resp,
	}
}

func (tracker *ProgressTracker) UpdateProgress(progress int32, des string) {
	info := newProgress(progress, des)
	tracker.msgCh <- info
}

func (tracker *ProgressTracker) getMsgChan() chan *proto.ActionProgress {
	return tracker.msgCh
}
