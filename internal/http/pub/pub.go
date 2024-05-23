package pub

import (
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/http/helper"
	"net/http"
	"os"
	"path/filepath"
)

// handleApiPubFile handles API management for public files.
func HandleApiUserPubFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		var File data.PubFile
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
		File.New(info.Name(), mime, info.Size(), username)
		if err := File.Save(); err != nil {
			http.Error(w, "Erro ao salvar arquivo", http.StatusInternalServerError)
			return
		}
		helper.JsonResponse(w, http.StatusCreated, helper.ResponseMap{"message": "Arquivo salvo com sucesso"})
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
