package http

import (
	"flag"
	"net/http"
)

var (
	rootDirectory = flag.String("root", "./root", "Diretório raiz do servidor WebDAV")
	globalToken   = flag.String("token", "123456", "Token de autenticação")
)

func InitMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/dav/", BasicAuthMiddleware(http.HandlerFunc(handleWebDAV)))
	mux.Handle("/admin/user", BearerAuthMiddleware(http.HandlerFunc(handleUserAdmin)))
	return mux
}
