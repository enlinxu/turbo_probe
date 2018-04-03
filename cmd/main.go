package main

import (
	"flag"
	"fmt"
	"github.com/golang/glog"
	"time"
	instance "turbo_probe/pkg/example"
	"turbo_probe/pkg/mediation"
	"turbo_probe/pkg/probe"
	"turbo_probe/pkg/restapi"
	"turbo_probe/pkg/wsocket"
)

const (
	defaultRemoteMediationServerPath string = "/vmturbo/remoteMediation"
	defaultRemoteMediationServerUser string = "vmtRemoteMediation"
	defaultRemoteMediationServerPwd  string = "vmtRemoteMediation"
)

var (
	probeCategory = "Cloud Native"
	probeType     = "Kubernetes.mock"
	protocolVer   = "6.1.0-SNAPSHOT"
	serverHost    = "https://localhost:9400/"

	username = "administrator"
	passwd   = "a"
)

func setFlags() {
	flag.StringVar(&serverHost, "serverHost", serverHost, "host of OpsMgr")
	flag.StringVar(&probeType, "probeType", probeType, "type of this probe")
	flag.StringVar(&protocolVer, "protocolVersion", protocolVer, "OpsMgr protocol version")
	flag.StringVar(&username, "turboUser", username, "OpsMgr user name")
	flag.StringVar(&passwd, "turboPasswd", passwd, "OpsMgr user password")
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

func getMediationClient() (*mediation.MediationClient, error) {
	//1. get websocket config
	conf := getConnConfig()

	//2. get turbo probe
	//protocolVer := "6.1.0-SNAPSHOT"
	//probeType := "mock-probe"
	//probeCategory := "cloudNative"
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

func buildTarget(cate, ttype string) *restapi.Target {
	target := &restapi.Target{
		Category: cate,
		Type:     ttype,
	}

	builder := restapi.NewInputFieldsBuilder()
	builder.With("targetIdentifier", "myTargetId").
		With("username", "developer").
		With("password", "pass")

	target.InputFields = builder.Create()
	//target.IdentifyingFields = []string{"Address"}

	return target
}

func addTarget(cate, ttype string) {
	time.Sleep(time.Second * 10)
	//1. construct target
	target := buildTarget(cate, ttype)

	//2. construct a restapi client
	client, err := restapi.NewRestClient(serverHost, username, passwd)
	if err != nil {
		glog.Errorf("Failed to create restAPI client: %v", err)
		return
	}

	//3. add the target
	resp, err := client.AddTarget(target)
	if err != nil {
		glog.Errorf("Failed to add target: %v, %v", err, resp)
		return
	}

	glog.V(2).Infof("Add target succeded: %v", resp)
	return
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

	go addTarget(probeCategory, probeType)
	client.Start()

	return
}
