package http

import (
	"github.com/gabrielmoura/davServer/internal/http/admin"
	"github.com/gabrielmoura/davServer/internal/http/file"
	"github.com/gabrielmoura/davServer/internal/http/helper"
	"github.com/gabrielmoura/davServer/internal/http/pub"
	"net/http"
)

func InitMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(handleRoot))
	mux.Handle("/dav/", BasicAuthMiddleware(http.HandlerFunc(handleWebDAV)))
	mux.Handle("/pub/{name}", http.HandlerFunc(handlePubFile))

	mux.Handle("/admin/user", BearerGlobalAuthMiddleware(http.HandlerFunc(admin.HandleUserAdmin)))

	mux.Handle("/user/file", BearerAuthMiddleware(http.HandlerFunc(file.HandleApiFile)))
	mux.Handle("/user/pub", BearerAuthMiddleware(http.HandlerFunc(pub.HandleApiUserPubFile)))
	return mux
}
func handleRoot(w http.ResponseWriter, r *http.Request) {
	helper.JsonResponse(w, http.StatusOK, helper.ResponseMap{"message": "Olá, bem-vindo ao servidor DAV!"})
}
