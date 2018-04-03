package mediation

import (
	"turbo_probe/pkg/probe"
	"turbo_probe/pkg/proto"
	"turbo_probe/pkg/version"
	"turbo_probe/pkg/wsocket"

	"fmt"
	"github.com/golang/glog"
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"os"
	"os/signal"
	"syscall"
)

// MediationClient will glue websocket and turbo_probe
type MediationClient struct {
	wsConfig        *wsocket.ConnectionConfig
	wsconn          *wsocket.WSconnection
	wsRetryDuration time.Duration

	protocolVersion string

	probe *probe.TurboProbe

	shouldStop bool
}

func NewMediationClient(wsconf *wsocket.ConnectionConfig, probe *probe.TurboProbe) *MediationClient {
	return &MediationClient{
		wsConfig:        wsconf,
		wsRetryDuration: defaultConnectionRetryDuration,
		protocolVersion: probe.ProbeInfoProvider.GetProtocolVersion(),
		probe:           probe,
		shouldStop: false,
	}
}

func (m *MediationClient) Start() error {
	m.handleSignal()
	for {
		glog.V(2).Infof("Begin protocol hand shake ...")
		flag := m.ProtocolHandShake()
		if !flag {
			err := fmt.Errorf("MediationClient failed to do protocol hand shake, terminating.")
			glog.Errorf(err.Error())
			return err
		}

		glog.V(2).Infof("Begin to serve server requests ...")
		m.WaitServerRequests()

		if m.shouldStop {
			glog.V(1).Infof("Mediation Client is stopped.")
			break
		}

		du := m.wsRetryDuration
		glog.Errorf("websocket is closed. Will re-connect in %v seconds.", du.Seconds())
		time.Sleep(du)
	}

	return nil
}

func (m *MediationClient) Stop() {
	m.shouldStop = true
	if m.wsconn != nil {
		m.wsconn.Stop()
	}
}

func (m *MediationClient) handleSignal() {
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			glog.V(1).Info("Received SIGINT signal.")
		case syscall.SIGTERM:
			glog.V(1).Info("Reived SIGTERM signal, closed websocket")
		case syscall.SIGQUIT:
			glog.V(1).Info("Reived SIGQUIT signal, closed websocket")
		}
		m.Stop()
		time.Sleep(time.Second*10)
		glog.V(1).Infof("done xxx")
		time.Sleep(time.Second*2)
		os.Exit(0)
	}()
}


func (m *MediationClient) WaitServerRequests() {
	m.wsconn.Start()

	for {
		if m.shouldStop {
			glog.V(1).Info("Stop waiting for server request: websocket is closed.")
			return
		}

		//1. get request from server, and handle it
		datch, err := m.wsconn.GetReceived()
		if err != nil {
			glog.Errorf("Stop waiting for server request: %v", err)
			return
		}

		timer := time.NewTimer(time.Second * 10)
		select {
		case dat := <-datch:
			if m.shouldStop {
				glog.V(1).Info("Stop waiting for server request: websocket is closed.")
				return
			}
			go m.handleServerRequest(dat)
		case <-timer.C:
			continue
		}
	}
}

func (m *MediationClient) ProtocolHandShake() bool {

	for {
		glog.V(2).Infof("begin to connect to server, and do protocol hand shake.")
		m.buildWSConnection()

		glog.V(2).Infof("begin to do protocol hand shake")
		flag, err := m.doProtocolHandShake()
		if err == nil {
			return flag
		}

		if !flag {
			return false
		}

		du := time.Second * 20
		glog.Errorf("protocolHandShake failed, will retry in %v seconds", du.Seconds())
		time.Sleep(du)
	}
}

func (m *MediationClient) buildWSConnection() error {

	if m.wsconn != nil {
		m.wsconn.Stop()
		m.wsconn = nil
	}

	for {
		if m.shouldStop {
			return fmt.Errorf("Stopped")
		}

		wsconn := wsocket.NewConnection(m.wsConfig)
		if wsconn == nil {
			glog.Errorf("Failed to build websocket connection: %++v", m.wsConfig)
			glog.Errorf("Will Retry in %v seconds", m.wsRetryDuration.Seconds())
			time.Sleep(m.wsRetryDuration)
			continue
		}

		m.wsconn = wsconn
		break
	}

	return nil
}

func (m *MediationClient) negotiationVersion() (bool, error) {
	//1. negotiation protocol version
	request := &version.NegotiationRequest{
		ProtocolVersion: &m.protocolVersion,
	}

	dat_in, err := protobuf.Marshal(request)
	if err != nil {
		glog.Errorf("Failed to marshal Negotiation request(%++v): %v", request, err)
		return false, err
	}

	//2. send request and get answer
	dat_out, err := m.wsconn.SendRecv(dat_in, -1)
	if err != nil {
		glog.Errorf("Failed to get negotiation response: %v", err)
		// will retry
		return true, err
	}

	//3. parse the answer
	resp := &version.NegotiationAnswer{}
	if err := protobuf.Unmarshal(dat_out, resp); err != nil {
		glog.Errorf("Failed to unmarshal negotiaonAnswer(%s): %v", string(dat_out), err)
		//will retry
		return true, err
	}

	result := resp.GetNegotiationResult()
	if result != version.NegotiationAnswer_ACCEPTED {
		glog.Errorf("Protocol Version(%v) Negotiation is not accepted: %v", m.protocolVersion, resp.GetDescription())
		return false, nil
	}
	glog.V(2).Infof("Protocol Version Negotiaion success: %v", m.protocolVersion)

	return true, nil
}

func (m *MediationClient) registerProbe() (bool, error) {
	//1. probe info
	probeInfo, err := m.probe.ProbeInfoProvider.GetProbeInfo()
	if err != nil {
		glog.Errorf("Failed to get probeInfo: %v", err)
		return false, err
	}
	glog.V(3).Infof("Probe.Info=%++v", probeInfo)

	request := &proto.ContainerInfo{
		Probes: probeInfo,
	}

	dat_in, err := protobuf.Marshal(request)
	if err != nil {
		glog.Errorf("Failed to marshal probeInfo (%++v): %v", request, err)
		return false, err
	}

	//2. send request and get response
	dat_out, err := m.wsconn.SendRecv(dat_in, -1)
	if err != nil {
		glog.Errorf("Failed to get registration response: %v", err)
		return true, err
	}

	//3. parse the answer
	resp := &proto.Ack{}
	if err := protobuf.Unmarshal(dat_out, resp); err != nil {
		glog.Errorf("Failed to unmarshl registration ack(%s): %v", string(dat_out), err)
		return false, err
	}

	return true, nil
}

func (m *MediationClient) doProtocolHandShake() (bool, error) {

	//1. protocol version negotiation
	flag, err := m.negotiationVersion()
	if err != nil {
		glog.Errorf("protocolHandShake failed: %v", err)
		return flag, err
	}

	if !flag {
		glog.Errorf("protocolHandShake is not accepted: %s is not accepted", m.protocolVersion)
		return false, nil
	}
	glog.V(3).Infof("probe protocol version negotiation success")

	//2. register probe info
	flag, err = m.registerProbe()
	if err != nil {
		glog.Errorf("protocolHandShake failed: %v", err)
		return flag, err
	}
	glog.V(3).Infof("probe registration success")

	return true, nil
}
