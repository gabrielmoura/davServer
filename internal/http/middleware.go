package http

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gabrielmoura/davServer/config"
	"github.com/gabrielmoura/davServer/internal/data"
	"github.com/gabrielmoura/davServer/internal/log"
	"net/http"
	"strings"
)

const (
	AuthInvalid   = "Autenticação Inválida"
	FormatInvalid = "Formato de autenticação inválido"
	TokenInvalid  = "Token inválido"
	RequiredAuth  = "Autenticação necessária"
)

// BearerAuthMiddleware checks for a valid bearer token in the request.
func BearerAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, TokenInvalid, http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Bearer" {
			http.Error(w, FormatInvalid, http.StatusBadRequest)
			return
		}

		authData, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, FormatInvalid, http.StatusBadRequest)
			return
		}

		creds := strings.SplitN(string(authData), ":", 2)
		if len(creds) != 2 {
			http.Error(w, FormatInvalid, http.StatusBadRequest)
			return
		}

		username, password := creds[0], creds[1]
		user, err := data.GetUser(username)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("Usuário inválido: %s %s", r.RemoteAddr, username))
			http.Error(w, AuthInvalid, http.StatusUnauthorized)
			return
		}

		if data.GenerateMD5Hash(password) != user.Password {
			log.Logger.Error(fmt.Sprintf("Senha inválida: %s %s", r.RemoteAddr, password))
			http.Error(w, AuthInvalid, http.StatusUnauthorized)
			return
		}
		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// BearerGlobalAuthMiddleware checks for a valid bearer token in the request.
func BearerGlobalAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Authorization") != fmt.Sprintf("Bearer %s", config.Conf.GlobalToken) {
			log.Logger.Error(fmt.Sprintf("Global Token Inválido: %s - %s", r.RemoteAddr, r.Header.Get("Authorization")))
			http.Error(w, TokenInvalid, http.StatusUnauthorized)
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
			http.Error(w, RequiredAuth, http.StatusUnauthorized)
			return
		}

		authParts := strings.SplitN(authHeader, " ", 2)
		if len(authParts) != 2 || authParts[0] != "Basic" {
			http.Error(w, FormatInvalid, http.StatusBadRequest)
			return
		}

		authData, err := base64.StdEncoding.DecodeString(authParts[1])
		if err != nil {
			http.Error(w, FormatInvalid, http.StatusBadRequest)
			return
		}

		creds := strings.SplitN(string(authData), ":", 2)
		if len(creds) != 2 {
			http.Error(w, FormatInvalid, http.StatusBadRequest)
			return
		}

		username, password := creds[0], creds[1]
		user, err := data.GetUser(username)
		if err != nil {
			log.Logger.Error(fmt.Sprintf("Usuário inválido: %s %s", r.RemoteAddr, username))
			http.Error(w, AuthInvalid, http.StatusUnauthorized)
			return
		}

		if data.GenerateMD5Hash(password) != user.Password {
			log.Logger.Error(fmt.Sprintf("Senha inválida: %s %s", r.RemoteAddr, password))
			http.Error(w, AuthInvalid, http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), "user", user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
