// Package model provides the base data types used in the task system.
package model

type TodoItem struct {
	Id      uint64 `json:"id"`
	Summary string `json:"summary"`
	Done    bool   `json:"done"`
}
