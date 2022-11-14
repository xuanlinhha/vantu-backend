package model

import (
	"encoding/json"

	"go.uber.org/zap"
	"vantu.org/go-backend/common"
)

type Phrase struct {
	Id      string `json:"id"`
	Han     string `json:"han"`
	Content string `json:"content"`
	Info    string `json:"info"`
	// Svg     string `json:"svg"`
}

type PhraseJson struct {
	Id      string                 `json:"id"`
	Han     string                 `json:"han"`
	Content ContentJson            `json:"content"`
	Info    map[string]interface{} `json:"info"`
	// Svg     string                 `json:"svg"`
}

type ContentJson struct {
	NguyenDu  map[string]interface{} `json:"nguyendu,omitempty"`
	ThieuChuu map[string]interface{} `json:"thieuchuu,omitempty"`
	Vdict     map[string]interface{} `json:"vdict,omitempty"`
	Mdbg      map[string]interface{} `json:"mdbg,omitempty"`
	Arch      map[string]interface{} `json:"arch,omitempty"`
}

func ConvertToJson(p *Phrase) *PhraseJson {
	var content ContentJson
	if err := json.Unmarshal([]byte(p.Content), &content); err != nil {
		common.Logger.Error("p Content: ", zap.Error(err))
	}
	var info map[string]interface{}
	if err := json.Unmarshal([]byte(p.Info), &info); err != nil {
		common.Logger.Error("p Info: ", zap.Error(err))
	}
	// var svg string
	// if err := json.Unmarshal([]byte(p.Svg), &svg); err != nil {
	// 	common.Logger.Error("p Svg: ", zap.Error(err))
	// }
	pj := PhraseJson{
		Id:      p.Id,
		Han:     p.Han,
		Content: content,
		Info:    info,
		// Svg:     svg,
	}
	return &pj
}
