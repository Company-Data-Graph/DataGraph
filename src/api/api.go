package api

import (
	"encoding/json"
	"fmt"
	"log"
	"media-server/src/models"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type MediaAPI struct {
	Host     string
	Port     int
	RootPath string
	Routes   models.MediaAPIRoutes
}

func NewMediaAPI(config *models.MediaAPIConfig) (*MediaAPI, error) {
	api := &MediaAPI{Host: config.Host, Port: config.Port, RootPath: config.StorageRootPath, Routes: config.Routes}
	users = make(map[string]string)
	users["admin"] = config.AdminPass
	return api, nil
}

func (api *MediaAPI) setRapidHeader(r *http.Request) {
	r.Header.Add("content-type", "application/x-www-form-urlencoded")
	r.Header.Add("Accept-Encoding", "application/gzip")
	r.Header.Add("X-RapidAPI-Key", "1ae2aa72f1mshb729e169c8532bfp14534cjsn7ddb25386cc2")
	r.Header.Add("X-RapidAPI-Host", "google-translate1.p.rapidapi.com")

}

func (api *MediaAPI) authorization(r *http.Request) int {
	token := r.Header.Get("Token")
	if token == "" {
		return http.StatusUnauthorized
	}
	log.Println(token)
	var claims models.Claims
	tokenValidation, err := jwt.ParseWithClaims(token, &claims, func(t *jwt.Token) (interface{}, error) { return jwtKey, nil })
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			return http.StatusUnauthorized
		}
		log.Println(err)
		return http.StatusBadRequest
	}
	if !tokenValidation.Valid {
		return http.StatusUnauthorized
	}
	log.Println(claims.Username, " authorized suc!")
	return http.StatusOK
}

func (api *MediaAPI) SetCorsHeaders(rw *http.ResponseWriter) {
	(*rw).Header().Set("Access-Control-Allow-Origin", "*")
	(*rw).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (api *MediaAPI) Run() {
	http.HandleFunc(api.Routes.DataRoute.Name, api.getDataByUrl)
	http.HandleFunc("/signin", api.signIn)
	log.Printf("Run server on %s:%d", api.Host, api.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", api.Host, api.Port), nil)
}

func (api *MediaAPI) signIn(rw http.ResponseWriter, r *http.Request) {
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	if !CredentialsValidation(creds) {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	token, err := GenerateNewToken(creds, 2)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	fmt.Fprintf(rw, fmt.Sprintf("{\"token\": \"%s\"}", token))
}

func (api *MediaAPI) getDataByUrl(rw http.ResponseWriter, r *http.Request) {
	currentPath := r.URL.RequestURI()[len(api.Routes.DataRoute.Name):]
	http.ServeFile(rw, r, fmt.Sprintf("%s%s%s", api.RootPath, api.Routes.DataRoute.StorageRoute, currentPath))
}
