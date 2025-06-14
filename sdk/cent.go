package sdk

import (
	"encoding/json"
	"github.com/Baidu-AIP/golang-sdk/aip/censor"
)

type Cont struct {
	ConclusionType uint64 `json:"conclusionType"`
}

func Content(content string) bool {
	client := censor.NewCloudClient("AhHky3pM1jZa5QmCRGIL5VUS", "cQeYQc7XW2WAjLm45y4MgSMzktrIoE1y")

	res := client.TextCensor(content)

	var data Cont

	json.Unmarshal([]byte(res), &data)
	if data.ConclusionType == 1 {
		return true
	}
	return false
}
