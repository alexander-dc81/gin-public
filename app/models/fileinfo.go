package models

//FileInfo model entity
type FileInfo struct {
	ContentType string
	Filename    string
	RealFormat  string `json:",omitempty"`
	Size        int
	Status      string `json:",omitempty"`
}
