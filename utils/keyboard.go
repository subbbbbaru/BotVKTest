package utils

import "encoding/json"

type JSONObject interface {
	ToJSON() string
}

type Button struct {
	Action struct {
		Type    string `json:"type"`
		Payload string `json:"payload,omitempty"`
		Label   string `json:"label"`
	} `json:"action"`
	Color string `json:"color"`
}

type Keyboard struct {
	OneTime bool       `json:"one_time"`
	Inline  bool       `json:"inline"`
	Buttons [][]Button `json:"buttons"`
}

func NewButton(label string, payload interface{}) Button {
	button := Button{}
	button.Action.Type = "text"
	button.Action.Label = label
	button.Action.Payload = ""
	if payload != nil {
		jPayoad, err := json.Marshal(payload)
		if err == nil {
			button.Action.Payload = string(jPayoad)
		}
	}
	button.Color = "primary"
	return button
}

func (keyboard Keyboard) ToJSON() string {

	keyboardJSON, _ := json.Marshal(keyboard)
	return string(keyboardJSON)
}
