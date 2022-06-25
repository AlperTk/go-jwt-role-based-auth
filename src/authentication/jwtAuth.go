package authentication

import (
	authorization "JwtAuth/src/authorization/service"
	"JwtAuth/src/errors"
	"crypto/tls"
	"encoding/json"
	errors2 "errors"
	"github.com/gorilla/mux"
	"net/http"
	"strings"
)

type JwtAuth struct {
	TokenProcessor TokenProcessor
	RoleAuthor     authorization.Authorizer
}

func init() {
	cfg := &tls.Config{
		InsecureSkipVerify: true,
	}
	http.DefaultClient.Transport = &http.Transport{
		TLSClientConfig: cfg,
	}
}

func (j JwtAuth) SetupMux(router *mux.Router) {
	router.Use(j.protect)
	router.NotFoundHandler = j.notFoundHandler()
}

func (j JwtAuth) protect(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		roles, err := j.tokenValidate(r)
		if err != nil {
			responseAuthorized(w)
			return
		}

		if j.RoleAuthor != nil {
			_, e := j.RoleAuthor.Process(roles, w, r)
			if e != nil {
				return
			}
		}

		next.ServeHTTP(w, r)
	})
}

func (j JwtAuth) notFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Content-Type", "application/json")
		_, err := j.tokenValidate(r)

		if err != nil {
			responseAuthorized(w)
			return
		}

		notFound(w)
	})
}

func (j JwtAuth) tokenValidate(r *http.Request) ([]string, error) {
	authHeader := r.Header.Get("Authorization")

	if len(authHeader) < 1 {
		return nil, errors2.New("no authorization token find")
	}

	bearerToken := strings.Split(authHeader, " ")[1]

	valid, roles, err := j.TokenProcessor.Process(bearerToken, r)
	if err != nil {
		return nil, err
	}

	if !valid {
		return nil, errors2.New("token not valid")
	}
	return roles, nil
}

func responseAuthorized(w http.ResponseWriter) {
	w.WriteHeader(401)
	_ = json.NewEncoder(w).Encode(errors.UnauthorizedError())
}

func notFound(w http.ResponseWriter) {
	w.WriteHeader(404)
	_ = json.NewEncoder(w).Encode(errors.NotFound())
}