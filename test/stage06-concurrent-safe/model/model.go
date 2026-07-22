package model

type Server struct {
	Name   string `json:"name"`
	IP     string `json:"ip"`
	Status string `json:"status"`
}
