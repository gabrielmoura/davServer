package http

import (
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/pkg/ternary"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// handleUserAdmin handles user management requests.
func handleUserAdmin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		username := r.FormValue("username")
		password := r.FormValue("password")
		if username == "" || password == "" {
			http.Error(w, "Usuário e senha são obrigatórios", http.StatusBadRequest)
			return
		}

		if _, err := data.GetUser(username); err == nil {
			http.Error(w, "Usuário já existe", http.StatusConflict)
			return
		}

		userDir := filepath.Join(config.Conf.ShareRootDir, username)
		if _, err := os.Stat(userDir); os.IsNotExist(err) {
			if _, err := data.CreateUserDirectory(userDir); err != nil {
				http.Error(w, "Erro ao criar pasta do usuário", http.StatusInternalServerError)
				return
			}
		}

		data.CreateUser(data.User{
			Username: username,
			Password: data.GenerateMD5Hash(password),
		})

		jsonResponse(w, http.StatusCreated, data.ResponseMap{"message": "Usuário criado com sucesso"})

	case http.MethodGet:
		jsonResponse(w, http.StatusOK, data.ResponseMap{"users": data.GetValidUsers()})

	case http.MethodDelete:
		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "Usuário é obrigatório", http.StatusBadRequest)
			return
		}
		data.DeleteUser(username)
		jsonResponse(w, http.StatusNoContent, data.ResponseMap{"message": fmt.Sprintf("Usuário %s removido com sucesso", username)})

	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// handleWebDAV handles WebDAV requests.
func handleWebDAV(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(data.User)
	if !ok || user.Username == "" {
		http.Error(w, "Usuário inválido", http.StatusUnauthorized)
		return
	}

	userDir := filepath.Join(config.Conf.ShareRootDir, user.Username)
	if _, err := os.Stat(userDir); os.IsNotExist(err) {
		http.Error(w, "Pasta do usuário não encontrada", http.StatusNotFound)
		return
	}

	fs := &webdav.Handler{
		FileSystem: webdav.Dir(userDir),
		LockSystem: webdav.NewMemLS(),
		Prefix:     "/dav",
		Logger: func(request *http.Request, err error) {
			if err != nil {
				log.Printf("Erro: %s %s: %v\n", request.Method, request.URL.Path, err)
			}
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
		fileData, err := fileToBase64(fullPath)
		if err != nil {
			http.Error(w, "Erro ao codificar arquivo em base64", http.StatusInternalServerError)
			return
		}
		fileResponse(w, http.StatusOK, ternary.OrString(name, metaFile.Name), metaFile.Size, metaFile.Mime, []byte(fileData))
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// handleApiFile handles API requests for files.
func handleApiFile(w http.ResponseWriter, r *http.Request) {}

// handleApiPubFile handles API management for public files.
func handleApiPubFile(w http.ResponseWriter, r *http.Request) {
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
		mime, err := getFileMimeType(fullPath)
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
		jsonResponse(w, http.StatusCreated, data.ResponseMap{"message": "Arquivo salvo com sucesso"})
	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
