package main

import (
	"context"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"golang.org/x/net/webdav"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"
)

var (
	rootDirectory = flag.String("root", "./root", "Diretório raiz do servidor WebDAV")
	globalToken   = flag.String("token", "123456", "Token de autenticação")
	users         []User
)

// User represents a user with a username and password.
type User struct {
	Username string
	Password string
}

// ResponseMap is a generic map for JSON responses.
type ResponseMap map[string]interface{}

// getValidUsers returns the list of valid users.
func getValidUsers() []User {
	return users
}

// createUser adds a new user to the list.
func createUser(user User) []User {
	users = append(users, user)
	return users
}

// deleteUser removes a user from the list.
func deleteUser(username string) []User {
	for i, u := range users {
		if u.Username == username {
			users = append(users[:i], users[i+1:]...)
			break
		}
	}
	return users
}

// getUser retrieves a user by username.
func getUser(username string) (User, error) {
	for _, u := range users {
		if u.Username == username {
			return u, nil
		}
	}
	return User{}, fmt.Errorf("usuário não encontrado")
}

// generateMD5Hash generates an MD5 hash of a given password.
func generateMD5Hash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// createUserDirectory creates a directory for a user.
func createUserDirectory(path string) (string, error) {
	if err := os.Mkdir(path, 0755); err != nil {
		return "", err
	}
	return path, nil
}

// BearerAuthMiddleware checks for a valid bearer token in the request.
func BearerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != *globalToken {
			http.Error(w, "Token inválido", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// BasicAuthMiddleware checks for valid basic auth credentials in the request.
func BasicAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.Header().Set("WWW-Authenticate", "Basic realm='WebDAV'")
			http.Error(w, "Autenticação necessária", http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, "Formato de autenticação inválido", http.StatusBadRequest)
			return
		}

		authData, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, "Formato de autenticação inválido", http.StatusBadRequest)
			return
		}

		creds := strings.SplitN(string(authData), ":", 2)
		if len(creds) != 2 {
			http.Error(w, "Formato de autenticação inválido", http.StatusBadRequest)
			return
		}

		username, password := creds[0], creds[1]
		user, err := getUser(username)
		if err != nil {
			http.Error(w, "Usuário inválido", http.StatusUnauthorized)
			return
		}

		if generateMD5Hash(password) != user.Password {
			http.Error(w, "Senha inválida", http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

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

		if _, err := getUser(username); err == nil {
			http.Error(w, "Usuário já existe", http.StatusConflict)
			return
		}

		userDir := filepath.Join(*rootDirectory, username)
		if _, err := os.Stat(userDir); os.IsNotExist(err) {
			if _, err := createUserDirectory(userDir); err != nil {
				http.Error(w, "Erro ao criar pasta do usuário", http.StatusInternalServerError)
				return
			}
		}

		createUser(User{
			Username: username,
			Password: generateMD5Hash(password),
		})

		w.WriteHeader(http.StatusCreated)
		result := ResponseMap{"message": "Usuário criado com sucesso"}
		w.Write([]byte(fmt.Sprintf("%v", result)))

	case http.MethodGet:
		result := ResponseMap{"users": getValidUsers()}
		w.Write([]byte(fmt.Sprintf("%v", result)))

	default:
		http.Error(w, "Método não permitido", http.StatusMethodNotAllowed)
	}
}

// handleWebDAV handles WebDAV requests.
func handleWebDAV(w http.ResponseWriter, r *http.Request) {
	user, ok := r.Context().Value("user").(User)
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

func main() {
	flag.Parse()

	mux := http.NewServeMux()
	mux.Handle("/dav/", BasicAuthMiddleware(http.HandlerFunc(handleWebDAV)))
	mux.Handle("/admin/user", BearerAuthMiddleware(http.HandlerFunc(handleUserAdmin)))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Channel to listen for termination signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Start server in a goroutine
	go func() {
		fmt.Println("Servidor WebDAV iniciado em dav://localhost:8080/dav/")
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Erro ao iniciar o servidor: %v\n", err)
		}
	}()

	// Block until a signal is received
	<-stop

	fmt.Println("\nIniciando o desligamento gracioso do servidor...")

	// Create a deadline for the shutdown
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Shutdown the server gracefully
	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Erro ao desligar o servidor: %v\n", err)
	}

	fmt.Println("Servidor desligado com sucesso.")
}
