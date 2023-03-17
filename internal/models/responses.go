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

type Error struct {
	What string `json:"error"`
}

type FileAlreadyExistError struct {
	What     string `json:"error"`
	FileName string `json:"fileName"`
}

type Token struct {
	Token string `json:"token"`
}

type FileName struct {
	FileName string `json: "fileName"`
}
