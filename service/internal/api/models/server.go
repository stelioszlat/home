package models

type Server struct {
	Port    int      `json:"port"`
	Host    string   `json:"host"`
	Modules []string `json:"modules"`
}
