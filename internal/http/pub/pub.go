package pub

import (
	"errors"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/http/helper"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// handleApiPubFile handles API management for public files.
func HandleApiUserPubFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		username := r.Context().Value("user").(data.User).Username
		// Pegue o caminho do arquivo
		path := r.FormValue("path")
		fullPath := filepath.Join(config.Conf.ShareRootDir, username, path)
		// verifica se o arquivo existe
		info, err := os.Stat(fullPath)
		if os.IsNotExist(err) {
			http.Error(w, "Arquivo não encontrado", http.StatusNotFound)
			return
		}
		mime, err := helper.GetFileMimeType(fullPath)
		if err != nil {
			http.Error(w, "Erro ao detectar tipo MIME", http.StatusInternalServerError)
			return
		}

		// Pegue informações do arquivo
		File := data.New(info.Name(), mime, info.Size(), username)
		if err := File.Save(); err != nil {
			if errors.Is(data.ErrPubArchiveExist, err) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			log.Printf("Error: ao salvar dado %s", err)
			http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
			return
		}
		helper.JsonResponse(w, http.StatusCreated, helper.ResponseMap{"message": "Arquivo salvo com sucesso", "file": File})
	case http.MethodGet:
		username := r.Context().Value("user").(data.User).Username
		pubFiles, err := data.ListPubFile(username)
		if err != nil {
			http.Error(w, "Erro", http.StatusServiceUnavailable)
		}
		helper.JsonResponse(w, http.StatusOK, pubFiles)
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
