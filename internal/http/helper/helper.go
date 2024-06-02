package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"io"
	"net/http"
	"os"
)

type ResponseMap map[string]interface{}

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
	if size > config.Conf.Srv.ChunkSize*1024*1024 {
		sendFileInChunks(w, file, size)
	} else {
		sendFile(w, code, name, size, mime, file)
	}
}

// sendFileInChunks envia o arquivo em chunks.
func sendFileInChunks(w http.ResponseWriter, file []byte, size int64) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Flusher not supported", http.StatusInternalServerError)
		return
	}

	buffer := make([]byte, 1024)
	reader := bytes.NewReader(file)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			break // Fim do arquivo
		}
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao ler arquivo: %s", err), http.StatusInternalServerError)
			return
		}

		_, err = w.Write(buffer[:n])
		if err != nil {
			http.Error(w, fmt.Sprintf("Erro ao escrever chunk: %s", err), http.StatusInternalServerError)
			return
		}

		flusher.Flush()
		w.Header().Set("Content-Range", fmt.Sprintf("bytes %d-%d/%d", reader.Size()-int64(reader.Len()), reader.Size()-int64(reader.Len())+int64(n)-1, size))
	}
}

// sendFile envia o arquivo inteiro de uma vez.
func sendFile(w http.ResponseWriter, code int, name string, size int64, mime string, file []byte) {
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
