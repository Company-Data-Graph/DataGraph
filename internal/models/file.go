package models

type FileModel struct {
	FileName   string `json:"fileName"`
	FileBuffer string `json:"fileBuffer"`
}
