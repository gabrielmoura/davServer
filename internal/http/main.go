package http

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

func InitMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/dav/", BasicAuthMiddleware(http.HandlerFunc(handleWebDAV)))
	mux.Handle("/pub/{name}", http.HandlerFunc(handlePubFile))
	mux.Handle("/admin/user", BearerGlobalAuthMiddleware(http.HandlerFunc(handleUserAdmin)))
	mux.Handle("/user/file", BearerAuthMiddleware(http.HandlerFunc(handleApiFile)))
	mux.Handle("/user/pub", BearerAuthMiddleware(http.HandlerFunc(handleApiPubFile)))
	return mux
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
