package login

import (
	"Orca_Master/define/retcode"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type userData struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type retData struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func UserLogin(url, username, password string) retData {
	payload, _ := json.Marshal(userData{Username: username, Password: password})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return retData{retcode.FAIL, "Failed to connect to server", ""}
	}
	body, _ := ioutil.ReadAll(resp.Body)
	var respData retData
	json.Unmarshal(body, &respData)
	defer resp.Body.Close()
	return respData
}
