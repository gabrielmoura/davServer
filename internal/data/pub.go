package data

import (
	"fmt"
	"github.com/google/uuid"
)

type PubFile struct {
	Name  string `json:"name"`  // file.png
	Hash  string `json:"hash"`  // 1234567890abcdef
	Mime  string `json:"mime"`  // image/png
	Size  int64  `json:"size"`  // 12345
	Owner string `json:"owner"` // Username
}

func (f PubFile) New(name, mime string, size int64, owner string) {
	f.Name = name
	f.Hash = uuid.New().String()
	f.Mime = mime
	f.Size = size
	f.Owner = owner
}
func (f PubFile) Save() error {
	return CreatePubFile(f)
}
func GetPubFile(hash string) (PubFile, error) {
	files, err := readPublicFiles(dB)
	if err != nil {
		return PubFile{}, err
	}
	for _, f := range files {
		if f.Hash == hash {
			return f, nil
		}
	}
	return PubFile{}, fmt.Errorf("arquivo não encontrado")
}
func CreatePubFile(file PubFile) error {
	files, err := readPublicFiles(dB)
	if err != nil {
		return err
	}
	files = append(files, file)
	return writePublicFiles(dB, files)
}
func DeletePubFile(hash string) error {
	files, err := readPublicFiles(dB)
	if err != nil {
		return err
	}
	for i, f := range files {
		if f.Hash == hash {
			files = append(files[:i], files[i+1:]...)
			break
		}
	}
	return writePublicFiles(dB, files)
}
