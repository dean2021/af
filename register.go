// Copyright 2020 Dean.
// Authors: Dean <dean@csoio.com>
// Date: 2020/5/20 4:00 下午

// agent注册相关

package af

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/dean2021/af/util/retry"
	"github.com/dean2021/af/util/system"
	"github.com/BurntSushi/toml"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

// 创建并保存agent注册信息到文件中
func createAgentRegisterInfoFile(tomlFile string, agent *Agent) error {
	f, err := os.OpenFile(tomlFile, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	encoder := toml.NewEncoder(f)
	err = encoder.Encode(agent)
	return err
}

// 注册agent
func register(agent *Agent) error {

	type RegisterResult struct {
		Id       string `json:"id"`
		ServerId int    `json:"serverId"`
		HostName string `json:"hostname"`
	}

	reqBody := map[string]interface{}{
		"hostName":  system.GetHostName(),
		"ipAddress": system.GetIPs(),
		"mac":       system.GetMacAddrs(),
	}

	jsonByte, err := json.Marshal(reqBody)
	if err != nil {
		log.Fatal(err)
	}

	// TODO 重写代理,用于调试接口
	client := &http.Client{}

	req, err := http.NewRequest("POST", agent.Config.Get("system.register.api"), bytes.NewBuffer(jsonByte))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("接口返回状态码非200: " + resp.Status)
	}

	if resp.Body == nil {
		return errors.New("接口返回内容为空")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	serverResp := &struct {
		Success bool           `json:"success"`
		Code    int            `json:"code"`
		Message string         `json:"message"`
		Result  RegisterResult `json:"result"`
	}{}

	err = json.Unmarshal(body, serverResp)
	if err != nil {
		return errors.New("接口返回格式错误:" + err.Error())
	}

	if serverResp.Success == false || serverResp.Code != 200 {
		return errors.New(fmt.Sprintf("接口返回错误提示:%s", serverResp.Message))
	}

	agent.ID = serverResp.Result.Id
	return nil
}

// 注册agent
func Register(agentInfo *Agent) {
	agentInfoFilePath := agentInfo.Config.Get("system.register.save_file")
	_, err := toml.DecodeFile(agentInfoFilePath, agentInfo)
	if err != nil && !os.IsNotExist(err) {
		log.Fatal(err)
	}

	// 首次注册
	if err != nil && os.IsNotExist(err) {

		log.Println("开始注册agent...")
		// 向服务端注册agent
		err := retry.Do(
			func() error {
				return register(agentInfo)
			},
			retry.Attempts(3),
			retry.Delay(time.Second),
			retry.LastErrorOnly(false),
			retry.OnRetry(func(n uint, err error) {
				log.Printf("Registration failed#%d: %s\n", n, err)
			}),
		)
		if err != nil {
			log.Fatal(err)
		}

		// 注册成功后创建agent info文件
		err = createAgentRegisterInfoFile(agentInfoFilePath, agentInfo)
		if err != nil {
			log.Fatalf("agent注册文件创建失败:%v", err)
		}

		// TODO 注册成功,增加服务端通知
		log.Println("注册成功, AgentID:" + agentInfo.ID + ", 存放文件: " + agentInfoFilePath)
	}
}
