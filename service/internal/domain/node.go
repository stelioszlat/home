package domain

import "time"

type Node struct {
	Name      string    `json:"name"`
	Os        string    `json:"os"`
	NodeType  string    `json:"-"`
	Memory    string    `json:"-"`
	CPU       string    `json:"-"`
	Storage   string    `json:"-"`
	Network   string    `json:"-"`
	IsActive  string    `json:"-"`
	createdAt time.Time `json:"-"`
	updatedAt time.Time `json:"-"`
}
