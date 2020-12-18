package models

type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error string      `json:"error,omitempty"`
	Meta  interface{} `json:"meta,omitempty"`
}
