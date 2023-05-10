package utils

type GroupsLongPollServer struct {
	Response LPResponse `json:"response"`
	Error    LPError    `json:"error"`
}

type LPResponse struct {
	Server string `json:"Server"`
	Key    string `json:"Key"`
	Ts     string `json:"Ts"`
}
type LPError struct {
	ErrorCode int    `json:"error_code"`
	ErrorMsg  string `json:"error_msg"`
}
