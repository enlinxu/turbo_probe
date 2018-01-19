package restapi

import (
	"net/url"
	"crypto/tls"
	"net/http"
	"github.com/golang/glog"
	"fmt"
	"bytes"
	"encoding/json"
	"io/ioutil"
)

type TurboRestClient struct {
	client *http.Client

	host     string
	username string
	password string
}

func NewRestClient(host, uname, pass string) (*TurboRestClient, error) {
	//1. get http client
	client := &http.Client{
		Timeout: defaultTimeOut,
	}

	//2. check whether it is using ssl
	addr, err := url.Parse(host)
	if err != nil {
		glog.Errorf("Invalid url:%v, %v", host, err)
		return nil, err
	}
	if addr.Scheme == "https" {
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		client.Transport = tr
	}

	return &TurboRestClient{
		client:   client,
		host:     host,
		username: uname,
		password: pass,
	}, nil
}

func (c *TurboRestClient) AddTarget(target *Target) (string, error) {
	//1. gen request
	req, err := c.genAddTargetRequest(target)
	if err != nil {
		glog.Errorf("Failed to generate AddTargetRquest: %v", err)
		return "", err
	}

	//2. send request and receive result
	resp, err := c.client.Do(req)
	if err != nil {
		glog.Errorf("Failed to send request: %v", err)
		return "", err
	}

	//3. parse result
	defer resp.Body.Close()
	content, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		glog.Errorf("Failed to read response: %v", err)
		return "Failed to read respone", err
	}

	result := string(content)
	code := resp.StatusCode
	glog.V(4).Infof("Add target code=%d, content=%v", code, result)

	if code >= 200 && code < 300 {
		glog.V(3).Infof("Add target succeded: %v", result)
		return result, nil
	}

	if code >= 400 && code < 400 {
		glog.Errorf("Failed to add request: %v", result)
		err := fmt.Errorf("Bad Request")
		return result, err
	}

	if code >= 500 {
		err := fmt.Errorf("API server internal error")
		glog.Errorf("Failed to add request: %v", err)
		return result, err
	}

	glog.Errorf("Failed to add target, bad status code: %v, %v", code, result)
	return result, fmt.Errorf("Bad status code: %v", code)
}

func (c *TurboRestClient) genAddTargetRequest(target *Target) (*http.Request, error) {
	//0. data
	data, err := json.Marshal(target)
	if err != nil {
		glog.Errorf("failed to generate json: %v", err)
		return nil, err
	}

	glog.V(4).Infof("target json data: %++v", string(data))

	//1. a new http request
	urlStr := fmt.Sprintf("%s%s", c.host, API_PATH_TARGET)
	req, err := http.NewRequest("POST", urlStr, bytes.NewBuffer(data))
	if err != nil {
		glog.Errorf("Failed to generate a http.request: %v", err)
		return nil, err
	}

	//2. set queries
	q := req.URL.Query()
	req.URL.RawQuery = q.Encode()

	//3. set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	//4. add user/password
	req.SetBasicAuth(c.username, c.password)

	return req, nil
}
