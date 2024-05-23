package admin

import (
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/http/helper"
	"net/http"
	"os"
	"path/filepath"
)

// handleUserAdmin handles user management requests.
func HandleUserAdmin(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		createUser(w, r)
	case http.MethodGet:
		helper.JsonResponse(w, http.StatusOK, helper.ResponseMap{"users": data.GetValidUsers()})
	case http.MethodDelete:
		deleteUser(w, r)

	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}
func deleteUser(w http.ResponseWriter, r *http.Request) {
	username := r.FormValue("username")
	if username == "" {
		http.Error(w, "Usuário é obrigatório", http.StatusBadRequest)
		return
	}
	data.DeleteUser(username)
	helper.JsonResponse(w, http.StatusNoContent, helper.ResponseMap{"message": fmt.Sprintf("Usuário %s removido com sucesso", username)})
}
func createUser(w http.ResponseWriter, r *http.Request) {
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

	helper.JsonResponse(w, http.StatusCreated, helper.ResponseMap{"message": "Usuário criado com sucesso"})
}
