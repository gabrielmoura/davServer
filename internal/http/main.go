package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/emersion/go-webdav/caldav"
	"github.com/emersion/go-webdav/carddav"
	iCal "github.com/gabrielmoura/davServer/internal/caldav"
	ICard "github.com/gabrielmoura/davServer/internal/carddav"
	"io"
	"log"
	"net/http"
	"os"
)

func InitMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/dav/", BasicAuthMiddleware(http.HandlerFunc(handleWebDAV)))
	mux.Handle("/caldav/", BasicAuthMiddleware(http.HandlerFunc(handleCalDav)))
	mux.Handle("/carddav/", BasicAuthMiddleware(http.HandlerFunc(handleCardDav)))
	mux.Handle("/.well-known/", http.HandlerFunc(handleWellKnown))
	mux.Handle("/pub/{name}", http.HandlerFunc(handlePubFile))
	mux.Handle("/admin/user", BearerGlobalAuthMiddleware(http.HandlerFunc(handleUserAdmin)))
	mux.Handle("/user/file", BearerAuthMiddleware(http.HandlerFunc(handleApiFile)))
	mux.Handle("/user/pub", BearerAuthMiddleware(http.HandlerFunc(handleApiPubFile)))
	return mux
}
func handleWellKnown(w http.ResponseWriter, r *http.Request) {
	log.Printf("Well-known request: %s %s", r.Method, r.URL.Path)
	if r.URL.Path == "/.well-known/caldav" {
		http.Redirect(w, r, "/caldav", http.StatusFound)
	}
	if r.URL.Path == "/.well-known/carddav" {
		http.Redirect(w, r, "/carddav", http.StatusFound)
	}
}
func handleCalDav(w http.ResponseWriter, r *http.Request) {
	back := &iCal.Backend{}
	handle := &caldav.Handler{
		Backend: back,
		Prefix:  "/caldav",
	}
	handle.ServeHTTP(w, r)
}
func handleCardDav(w http.ResponseWriter, r *http.Request) {
	//back := ICard.CardBackend{
	//}
	handle := &carddav.Handler{
		Backend: &ICard.CardBackend{},
		Prefix:  "/carddav",
	}
	handle.ServeHTTP(w, r)
}

// jsonResponse sends a JSON response with the given status code and data.
func jsonResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar resposta: %s", err), http.StatusInternalServerError)
	}
}

// fileResponse sends a file response with the given status code, name, size, mime type, and file content.
func fileResponse(w http.ResponseWriter, code int, name string, size int64, mime string, file []byte) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	w.WriteHeader(code)
	_, err := w.Write(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao escrever arquivo: %s", err), http.StatusInternalServerError)
	}
}

func getFileMimeType(filePath string) (string, error) {
	// Abrir o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Ler os primeiros 512 bytes do arquivo para detectar o tipo MIME
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	// Detectar o tipo MIME baseado no conteúdo
	mimeType := http.DetectContentType(buffer)
	return mimeType, nil
}

func fileToBase64(filePath string) (string, error) {
	// Abrir o arquivo
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Ler o conteúdo do arquivo
	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	// Codificar o conteúdo em base64
	base64Encoded := base64.StdEncoding.EncodeToString(fileBytes)

	return base64Encoded, nil
}
