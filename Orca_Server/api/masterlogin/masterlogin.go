package masterlogin

import (
	"Orca_Server/api"
	"Orca_Server/define/retcode"
	"Orca_Server/servers"
	"Orca_Server/setting"
	"encoding/json"
	"net/http"
)

type Controller struct {
}

type inputData struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
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

	err = servers.MasterLogin(inputData.Username, inputData.Password)
	if err != nil {
		api.Render(w, retcode.FAIL, err.Error(), []string{})
		return
	}

	api.Render(w, retcode.SUCCESS, "login success", setting.CommonSetting.CryptoKey)
	return
}
