package stager

import (
	"bytes"
	"encoding/json"
	"net/http"
)

type systemIdData struct {
	SystemId string `json:"systemId" validate:"required"`
}

// 注册SystemId
func registerSystemId(url, systemId string) {
	payload, _ := json.Marshal(systemIdData{SystemId: systemId})
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
}
