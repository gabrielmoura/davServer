package http

import (
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/http/helper"
	"github.com/gabrielmoura/davServer/internal/log"
	"github.com/gabrielmoura/go/pkg/ternary"
	"golang.org/x/net/webdav"
	"net/http"
	"os"
	"path/filepath"
)

// handleWebDAV handles WebDAV requests.
func handleWebDAV(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(data.User)
	if !ok || user.Username == "" {
		http.Error(w, "Usuário inválido", http.StatusUnauthorized)
		return
	}

	userDir := filepath.Join(config.Conf.ShareRootDir, user.Username)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		log.Logger.Error(fmt.Sprintf("Pasta do usuário não encontrada: %s", userDir))
		http.Error(w, "Pasta do usuário não encontrada", http.StatusNotFound)
		return
	}

	fs := &webdav.Handler{
		FileSystem: webdav.Dir(userDir),
		LockSystem: webdav.NewMemLS(),
		Prefix:     "/dav",
		Logger: func(request *http.Request, err error) {
			if err != nil {
				log.Logger.Error(fmt.Sprintf("%s WebDAV %s %s: %v\n", request.RemoteAddr, request.Method, request.URL.Path, err))
			}
			log.Logger.Debug(fmt.Sprintf("%s WebDAV %s %s\n", request.RemoteAddr, request.Method, request.URL.Path))
		},
	}

	fs.ServeHTTP(w, r)
}

// handlePubFile handles get requests for public files.
func handlePubFile(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		hash := r.URL.Query().Get("hash")
		name := r.URL.Query().Get("name")
		if hash == "" {
			http.Error(w, "Hash é obrigatório", http.StatusBadRequest)
			return
		}
		metaFile, err := data.GetPubFile(hash)
		if err != nil {
			http.Error(w, "Arquivo não encontrado", http.StatusNotFound)
			return
		}
		fullPath := filepath.Join(config.Conf.ShareRootDir, metaFile.Owner, metaFile.Name)
		fileData, err := helper.FileToBase64(fullPath)
		if err != nil {
			http.Error(w, "Erro ao codificar arquivo em base64", http.StatusInternalServerError)
			return
		}
		helper.FileResponse(w, http.StatusOK, ternary.OrString(name, metaFile.Name), metaFile.Size, metaFile.Mime, []byte(fileData))
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
