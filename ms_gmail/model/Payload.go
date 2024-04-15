package model

type RequestPayload struct {
	ObjectModel string      `json:"object"`
	Data        interface{} `json:"data"`
}

type Response struct {
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
	Success bool        `json:"success"`
}

// Use to test MTProto API ping status
type MTprotoPayload struct {
	Message string `json:"message"`
}

