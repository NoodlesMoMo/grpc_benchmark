package common

import (
	"encoding/json"
	"math/rand"
	"strings"
	"time"
)

type DummyData struct {
	Code int    `json:"code"`
	Data string `json:"data"`
}

func GenerateDummyData() []byte {
	x := rand.New(rand.NewSource(time.Now().UnixNano())).Intn(100)
	data := DummyData{
		Code: x,
		Data: strings.Repeat("x", x),
	}

	dataJson, _ := json.Marshal(data)

	return dataJson
}
