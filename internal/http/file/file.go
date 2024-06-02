package file

import (
	"context"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/http/helper"
	"github.com/gabrielmoura/davServer/internal/log"
	"go.uber.org/zap"
	"net/http"
	"os"
	"path/filepath"
)

type File struct {
	Name    string `json:"name"`
	Size    int64  `json:"size"`
	Mime    string `json:"mime"`
	Content []byte `json:"content"`
}
type Directory struct {
	Name  string `json:"name"`
	Files []File `json:"files"`
}

func listDir(path string) (Directory, error) {
	var dir Directory
	dir.Name = filepath.Base(path)
	files, err := os.ReadDir(path)
	if err != nil {
		return dir, err
	}
	for _, file := range files {
		if file.IsDir() {
			continue
		}
		content, err := os.ReadFile(filepath.Join(path, file.Name()))
		if err != nil {
			return dir, err
		}
		info, _ := file.Info()
		dir.Files = append(dir.Files, File{
			Name:    file.Name(),
			Size:    info.Size(),
			Mime:    http.DetectContentType(content),
			Content: content,
		})
	}
	return dir, nil
}
func readFile(path string) ([]byte, error) {
	return os.ReadFile(path)
}

// handleApiFile handles API requests for files.
func HandleApiFile(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(data.User)
	if !ok {
		http.Error(w, "Usuário inválido", http.StatusUnauthorized)
		return
	}
	ctx := context.WithValue(r.Context(), "userPath", filepath.Join(config.Conf.ShareRootDir, user.Username))

	switch r.Method {
	case http.MethodGet:
		if r.URL.Path == "/user/file" {
			getAllFile(w, r.WithContext(ctx))
		} else if r.URL.Path == "/user/file/{fileId}" {
			getFileById(w, r.WithContext(ctx))
		} else if r.URL.Path == "/user/file/{fileId}/metadata" {
			getFileMetadata(w, r.WithContext(ctx))
		} else {
			http.Error(w, "Rota não encontrada", http.StatusNotFound)
		}
	case http.MethodPost:
		uploadFile(w, r.WithContext(ctx))
	case http.MethodDelete:
		deleteFile(w, r.WithContext(ctx))
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

func getAllFile(w http.ResponseWriter, r *http.Request) {
	userPath := r.Context().Value("userPath").(string)
	files, err := listDir(userPath)
	if err != nil {
		http.Error(w, "Erro ao listar diretório", http.StatusInternalServerError)
		return
	}

	helper.JsonResponse(w, http.StatusOK, files)
}
func getFileById(w http.ResponseWriter, r *http.Request) {
	userPath := r.Context().Value("userPath").(string)
	fileId := r.URL.Query().Get("fileId")
	file, err := readFile(filepath.Join(userPath, fileId))

	if err != nil {
		http.Error(w, "Erro ao ler arquivo", http.StatusInternalServerError)
		return
	}
	helper.FileResponse(w, http.StatusOK, fileId, int64(len(file)), http.DetectContentType(file), file)
}
func getFileMetadata(w http.ResponseWriter, r *http.Request) {}
func uploadFile(w http.ResponseWriter, r *http.Request) {
	userPath := r.Context().Value("userPath").(string)
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		http.Error(w, "Erro ao analisar formulário", http.StatusInternalServerError)
		return
	}
	file, handler, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Erro ao obter arquivo", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	fileBytes := make([]byte, handler.Size)
	_, err = file.Read(fileBytes)
	if err != nil {
		log.Logger.Error("Erro ao ler arquivo", zap.Error(err))
		http.Error(w, "Erro ao ler arquivo", http.StatusInternalServerError)
		return
	}

	err = os.WriteFile(filepath.Join(userPath, handler.Filename), fileBytes, 0644)
	if err != nil {
		http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
		return
	}
	helper.JsonResponse(w, http.StatusCreated, "Arquivo salvo com sucesso")
}
func deleteFile(w http.ResponseWriter, r *http.Request) {
	userPath := r.Context().Value("userPath").(string)
	fileId := r.URL.Query().Get("fileId")
	err := os.Remove(filepath.Join(userPath, fileId))
	if err != nil {
		http.Error(w, "Erro ao excluir arquivo", http.StatusInternalServerError)
		return
	}
	helper.JsonResponse(w, http.StatusOK, "Arquivo excluído com sucesso")
}
