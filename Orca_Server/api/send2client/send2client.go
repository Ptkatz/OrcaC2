package send2client

import (
	"Orca_Server/api"
	"Orca_Server/define/retcode"
	"Orca_Server/servers"
	"encoding/json"
	"net/http"
)

type Controller struct {
}

type inputData struct {
	ClientId   string `json:"clientId" validate:"required"`
	SendUserId string `json:"sendUserId"`
	Code       int    `json:"code"`
	Msg        string `json:"msg"`
	Data       string `json:"data"`
	MessageId  string `json:"messageId"`
}

func (c *Controller) Run(w http.ResponseWriter, r *http.Request) {
	var inputData inputData
	if err := json.NewDecoder(r.Body).Decode(&inputData); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := api.Validate(inputData)
	if err != nil {
		api.Render(w, retcode.FAIL, err.Error(), []string{})
		return
	}

	//发送信息
	servers.SendMessage2Client(inputData.ClientId, inputData.SendUserId, inputData.Code, inputData.Msg, &inputData.Data, inputData.MessageId)

	api.Render(w, retcode.SUCCESS, "success", map[string]string{
		"messageId": inputData.MessageId,
	})
	return
}
