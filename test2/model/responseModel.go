package model

type Response struct {
	Message string      `json:"message,omitempty"`
	Result  interface{} `json:"result"`
}
