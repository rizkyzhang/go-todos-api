package models

type Todo struct {
	Id uint `json:"id"`
	Todo *string `json:"todo,omitempty"`
	IsCompleted *bool `json:"is_completed,omitempty"`
}