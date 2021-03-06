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
	"github.com/BurntSushi/toml"
	"github.com/dean2021/af/util/retry"
	"github.com/dean2021/af/util/system"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

// 注册agent
func doRegister(agent *Agent) error {
	systemInfo, err := json.Marshal(system.GetSystemInfo())
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", agent.Config.Get("system.register.api"), bytes.NewBuffer(systemInfo))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	if resp.Body == nil {
		return errors.New("接口返回内容为空")
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return errors.New("接口返回状态码非200: " + resp.Status)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	result := &struct {
		Success bool   `json:"success"`
		Code    int    `json:"code"`
		Message string `json:"message"`
		Result  struct {
			Id       string `json:"id"`
			ServerId int    `json:"serverId"`
			HostName string `json:"hostname"`
		} `json:"result"`
	}{}
	err = json.Unmarshal(body, result)
	if err != nil {
		return fmt.Errorf("接口返回格式错误: %v", err)
	}
	if result.Success == false || result.Code != 200 {
		return fmt.Errorf("接口返回错误提示:%s", result.Message)
	}
	agent.ID = result.Result.Id
	return nil
}

// 注册agent
func Register(agent *Agent) error {
	agentFilePath, err := filepath.Abs(agent.Config.Get("system.register.save_file"))
	if err != nil {
		return err
	}
	err = os.MkdirAll(filepath.Dir(agentFilePath), os.ModePerm)
	if err != nil {
		return err
	}
	// 向服务端注册agent
	err = retry.Do(
		func() error {
			return doRegister(agent)
		},
		retry.Attempts(3),
		retry.Delay(time.Second),
		retry.LastErrorOnly(false),
		retry.OnRetry(func(n uint, err error) {
			agent.logger.Printf("Registration failed#%d: %s\n", n, err)
		}),
	)
	if err != nil {
		return err
	}
	// 注册成功后创建agent info文件
	f, err := os.OpenFile(agentFilePath, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	encoder := toml.NewEncoder(f)
	return encoder.Encode(agent)
}
