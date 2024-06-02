package data

import (
	"errors"
	"github.com/dgraph-io/badger/v4"
	"github.com/google/uuid"
	"log"
)

var (
	ErrPubArchiveExist = errors.New("este arquivo já está público")
	ErrArchiveNotFound = errors.New("arquivo não encontrado")
)

type PubFile struct {
	Name  string `json:"name"`  // file.png
	Hash  string `json:"hash"`  // 1234567890abcdef
	Mime  string `json:"mime"`  // image/png
	Size  int64  `json:"size"`  // 12345
	Owner string `json:"owner"` // Username
}

func New(name, mime string, size int64, owner string) *PubFile {
	return &PubFile{
		Hash:  uuid.New().String(),
		Name:  name,
		Mime:  mime,
		Size:  size,
		Owner: owner,
	}
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
	return PubFile{}, ErrArchiveNotFound
}

func CreatePubFile(file PubFile) error {
	files, err := readPublicFiles(dB)
	if err != nil {
		if errors.Is(badger.ErrKeyNotFound, err) {
			fileSlice := make([]PubFile, 0)
			fileSlice = append(fileSlice, file)
			return writePublicFiles(dB, fileSlice)
		}
		log.Printf("Erro: Erro ao buscar matriz %s", err)
		return err
	}
	if isRepeated(files, file) {
		return ErrPubArchiveExist
	}
	files = append(files, file)
	return writePublicFiles(dB, files)
}
func isRepeated(slice []PubFile, file PubFile) bool {
	for _, f := range slice {
		if f.Name == file.Name {
			return true
		}
	}
	return false
}
func ListPubFile(owner string) ([]PubFile, error) {
	files := make([]PubFile, 0)
	allFiles, err := readPublicFiles(dB)
	if err != nil {
		return files, err
	}
	for _, file := range allFiles {
		if file.Owner == owner {
			files = append(files, file)
		}
	}
	return files, nil

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
