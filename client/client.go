package client

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"time"
)

type JmsClient struct {
	*http.Client
	serverUrl string
}

func NewJmsClient(serverUrl, accessKey, secretKey string, timeout int) *JmsClient {
	return &JmsClient{
		Client: &http.Client{
			Transport: NewSigAuthRoundTripper(accessKey, secretKey),
			Timeout:   time.Duration(timeout) * time.Second,
		},
		serverUrl: serverUrl,
	}
}

// 根据 IP 查询资产 ID, 支持的协议
func (c *JmsClient) QueryIDByIP(ip string) (id string, protocols []string) {
	endpoint, _ := url.JoinPath(c.serverUrl, "api/v1/assets/assets/suggestions/")
	u, _ := url.Parse(endpoint)
	u.RawQuery = url.Values{"address": []string{ip}}.Encode()
	req, _ := http.NewRequest("GET", u.String(), nil)
	resp, err := c.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data []map[string]any
	_ = json.Unmarshal(body, &data)
	if len(data) == 0 {
		log.Fatal("❌查无资产")
	}
	id = data[0]["id"].(string)
	// 取出支持的协议列表
	for _, v := range data[0]["protocols"].([]any) {
		protocols = append(protocols, v.(map[string]any)["name"].(string))
	}
	return
}

// 下载用于连接资产的令牌
func (c *JmsClient) GenRDPToken(asset, account string) (string, error) {
	endpoint, _ := url.JoinPath(c.serverUrl, "api/v1/authentication/connection-token/")
	reqData := map[string]string{
		"account":        account,
		"asset":          asset,
		"connect_method": "mstsc",
		"protocol":       "rdp",
	}
	payload, _ := json.Marshal(reqData)
	req, _ := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	var data map[string]any
	_ = json.Unmarshal(body, &data)
	return data["id"].(string), nil
}

// 下载令牌对应的 RDP 文件
func (c *JmsClient) DownRDP(token, fullscreen string) (string, error) {
	endpoint, _ := url.JoinPath(c.serverUrl, "api/v1/authentication/connection-token", token, "rdp-file/")
	req, _ := http.NewRequest("GET", endpoint+"?full_screen="+fullscreen, nil)
	resp, err := c.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode == 200 {
		f, _ := os.CreateTemp("", "*.rdp")
		defer f.Close()
		f.Write(body)
		return f.Name(), nil
	} else {
		var data map[string]any
		_ = json.Unmarshal(body, &data)
		return "", errors.New(data["detail"].(string))
	}
}
