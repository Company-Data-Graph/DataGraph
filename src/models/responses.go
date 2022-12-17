package models

import "time"

type Extension struct {
	Name  string `json:"name"`
	IsDir bool   `json:"isDir"`
}

type File struct {
	Name    string    `json:"name"`
	ModTime time.Time `json:"date"`
}
