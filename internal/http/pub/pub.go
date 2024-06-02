package pub

import (
	"errors"
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/http/helper"
	"github.com/gabrielmoura/davServer/internal/log"
	"github.com/gabrielmoura/davServer/internal/msg"
	"go.uber.org/zap"
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
			http.Error(w, msg.FileNotFound, http.StatusNotFound)
			return
		}
		mime, err := helper.GetFileMimeType(fullPath)
		if err != nil {
			http.Error(w, msg.ErrDetectingMIME, http.StatusInternalServerError)
			return
		}

		// Pegue informações do arquivo
		File := data.New(info.Name(), mime, info.Size(), username)
		if err := File.Save(); err != nil {
			if errors.Is(data.ErrPubArchiveExist, err) {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
			log.Logger.Error(fmt.Sprintf("%s %s", msg.ErrSavingFile, info.Name()), zap.Error(err))
			http.Error(w, msg.ErrSavingFile, http.StatusInternalServerError)
			return
		}
		helper.JsonResponse(w, http.StatusCreated, helper.ResponseMap{"message": msg.SuccFileSaved, "file": File})
	case http.MethodGet:
		username := r.Context().Value("user").(data.User).Username
		pubFiles, err := data.ListPubFile(username)
		if err != nil {
			http.Error(w, "Erro", http.StatusServiceUnavailable)
		}
		helper.JsonResponse(w, http.StatusOK, pubFiles)
	case http.MethodDelete:
		hash := r.FormValue("hash")
		err := data.DeletePubFile(hash)
		if err != nil {
			http.Error(w, msg.ErrDelete, http.StatusInternalServerError)
			return
		}
		helper.JsonResponse(w, http.StatusOK, helper.ResponseMap{"message": msg.SuccRemoved})
	default:
		http.Error(w, msg.MethodNotAllowed, http.StatusMethodNotAllowed)
	}
}
