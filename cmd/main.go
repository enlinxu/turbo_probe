package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"turbo_probe/pkg/wsocket"
	instance "turbo_probe/pkg/example"
	"turbo_probe/pkg/probe"
	"turbo_probe/pkg/mediation"
)

const (
	defaultRemoteMediationServerPath       string = "/vmturbo/remoteMediation"
	defaultRemoteMediationServerUser   string = "vmtRemoteMediation"
	defaultRemoteMediationServerPwd    string = "vmtRemoteMediation"
)

var (
	serverHost = "https://localhost:9400/"
)

func setFlags() {
	flag.StringVar(&serverHost, "serverHost", serverHost, "host of OpsMgr")
}

func getConnConfig() *wsocket.ConnectionConfig {
	path := defaultRemoteMediationServerPath
	user := defaultRemoteMediationServerUser
	passwd := defaultRemoteMediationServerPwd

	conf, err := wsocket.NewConnectionConfig(serverHost, path, user, passwd)
	if err != nil {
		glog.Errorf("Failed to create connection config: %v", err)
		return nil
	}
	return conf
}

func getMediationClient() (*mediation.MediationClient, error){
	//1. get websocket config
	conf := getConnConfig()

	//2. get turbo probe
	protocolVer := "6.1.0-SNAPSHOT"
	probeType := "mock-probe"
	probeCategory := "cloudNative"
	infoProvider := instance.NewMockProbeInfoProvider(protocolVer, probeType, probeCategory)
	discoveryExecutor := instance.NewMockDiscoveryExecutor("mocke discovery executor")
	actionExecutor := instance.NewMockActionExecutor("mock action executor")

	probeBuilder := probe.NewTurboProbeBuilder()
	probeBuilder.WithRegInfoProvider(infoProvider)
	probeBuilder.WithDiscoveryExecutor(discoveryExecutor)
	probeBuilder.WithActionExecutor(actionExecutor)
	myprobe, err := probeBuilder.Create()
	if err != nil {
		glog.Errorf("Failed to create turbo Probe: %v", err)
		return nil, fmt.Errorf("Failed to create turboProbe")
	}

	mclient := mediation.NewMediationClient(conf, myprobe)
	return mclient, nil
}

func main() {
	setFlags()
	flag.Parse()

	fmt.Println("Start remote mediation client")
	glog.V(1).Infof("Start remote mediation client")

	client, err := getMediationClient()
	if err != nil {
		glog.Errorf("Failed to create mediation client: %v", err)
		return
	}

	client.Start()

	return
}
