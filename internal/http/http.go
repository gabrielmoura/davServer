package http

import (
	"fmt"
	"github.com/gabrielmoura/davServer/internal/data"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

// handleUserAdmin handles user management requests.
func handleUserAdmin(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Authorization") != *globalToken {
		http.Error(w, "Token inválido", http.StatusUnauthorized)
		return
	}

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

		userDir := filepath.Join(*rootDirectory, username)
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

		w.WriteHeader(http.StatusCreated)
		result := data.ResponseMap{"message": "Usuário criado com sucesso"}
		w.Write([]byte(fmt.Sprintf("%v", result)))

	case http.MethodGet:
		result := data.ResponseMap{"users": data.GetValidUsers()}
		w.Write([]byte(fmt.Sprintf("%v", result)))

	case http.MethodDelete:
		username := r.FormValue("username")
		if username == "" {
			http.Error(w, "Usuário é obrigatório", http.StatusBadRequest)
			return
		}
		data.DeleteUser(username)
		w.WriteHeader(http.StatusNoContent)
		w.Write([]byte(fmt.Sprintf("Usuário %s removido com sucesso", username)))

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

	userDir := filepath.Join(*rootDirectory, user.Username)
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
