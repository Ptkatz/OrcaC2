package servers

import (
	"Orca_Server/sqlmgmt"
	"errors"
	"sync"
	"time"
)

type accountInfo struct {
	SystemId     string `json:"systemId"`
	RegisterTime int64  `json:"registerTime"`
}

var SystemMap sync.Map

func Register(systemId string) (err error) {
	//校验是否为空
	if len(systemId) == 0 {
		return errors.New("系统ID不能为空")
	}

	accountInfo := accountInfo{
		SystemId:     systemId,
		RegisterTime: time.Now().Unix(),
	}

	if _, ok := SystemMap.Load(systemId); ok {
		return errors.New("该系统ID已被注册")
	}

	SystemMap.Store(systemId, accountInfo)

	return nil
}

func MasterLogin(username, password string) (err error) {
	//校验是否为空
	if len(username) == 0 || len(password) == 0 {
		return errors.New("用户名或密码不能为空")
	}
	//校验用户名和密码是否正确
	user := sqlmgmt.Login(username, password)
	if user.ID == 0 {
		return errors.New("用户名或密码不正确")
	}
	return nil
}
