package helper

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/vmihailenco/msgpack/v5"
	"io"
	"net/http"
	"os"
)

type ResponseMap map[string]interface{}

func MsgResponse(w http.ResponseWriter, code int, data interface{}) {
	b, err := msgpack.Marshal(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar resposta: %s", err), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(code)
	w.Write(b)
}

// jsonResponse sends a JSON response with the given status code and data.
func JsonResponse(w http.ResponseWriter, code int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao codificar resposta: %s", err), http.StatusInternalServerError)
	}
}

// fileResponse sends a file response with the given status code, name, size, mime type, and file content.
func FileResponse(w http.ResponseWriter, code int, name string, size int64, mime string, file []byte) {
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s", name))
	w.Header().Set("Content-Type", mime)
	w.Header().Set("Content-Length", fmt.Sprintf("%d", size))
	w.WriteHeader(code)
	_, err := w.Write(file)
	if err != nil {
		http.Error(w, fmt.Sprintf("Erro ao escrever arquivo: %s", err), http.StatusInternalServerError)
	}
}

func GetFileMimeType(filePath string) (string, error) {
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

func FileToBase64(filePath string) (string, error) {
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
