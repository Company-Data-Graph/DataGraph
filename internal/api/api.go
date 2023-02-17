package api

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"media-server/internal/domain"
	"media-server/internal/jwt"
	"net/http"
	"os"
	"regexp"
	"strings"

	jwtGo "github.com/dgrijalva/jwt-go"
)

type MediaAPI struct {
	Host             string
	Port             int
	Prefix           string
	TokenLiveTime    int
	RootPath         string
	DataStorageRoute string
}

func NewMediaAPI(config *domain.MediaAPIConfig) (*MediaAPI, error) {
	api := &MediaAPI{Host: config.Host, Port: config.Port, Prefix: config.Prefix, TokenLiveTime: config.TokenLiveTime, RootPath: config.StorageRootPath, DataStorageRoute: config.DataStorageRoute}
	jwt.Users = make(map[string]string)
	jwt.Users["admin"] = config.AdminPass
	return api, nil
}

func (api *MediaAPI) getFileExtension(fileName string) string {
	fileExtension := "_"
	fileNameSplitted := (strings.Split(fileName, "."))
	if len(fileNameSplitted) > 0 {
		fileExtension = fileNameSplitted[len(fileNameSplitted)-1]
	}
	return fileExtension
}

func (api *MediaAPI) encodeFileName(fileName string, fileExtension string) string {
	encodedFileName := md5.Sum([]byte(fileName))
	return fmt.Sprintf("%s.%s", hex.EncodeToString(encodedFileName[:]), fileExtension)
}

func (api *MediaAPI) getFullFilePath(fileExtension string) string {
	path := fmt.Sprintf("%s/%s/%s", api.RootPath, api.DataStorageRoute, fileExtension)
	reg := regexp.MustCompile("(/)*")
	return reg.ReplaceAllString(path, "$1")
}

func (api *MediaAPI) authorization(token string) int {
	var claims domain.Claims
	tokenValidation, err := jwtGo.ParseWithClaims(token, &claims, func(t *jwtGo.Token) (interface{}, error) { return jwt.Key, nil })
	if err != nil {
		if err == jwtGo.ErrSignatureInvalid {
			return http.StatusUnauthorized
		}
		return http.StatusUnauthorized
	}
	if !tokenValidation.Valid {
		return http.StatusUnauthorized
	}
	return http.StatusOK
}

func (api *MediaAPI) setCorsHeaders(rw *http.ResponseWriter) {
	(*rw).Header().Set("Access-Control-Allow-Origin", "*")
	(*rw).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func (api *MediaAPI) checkMethod(getted string, required string) bool {
	if required != getted {
		return false
	}
	return true
}

func (api *MediaAPI) Run() {
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/ping/"), api.ping)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/data/"), api.getDataByUrl)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/signin"), api.signIn)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/upload"), api.upload)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/delete"), api.delete)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/dir"), api.getFileNamesWithDates)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/extensions"), api.getAvailableExtension)
	log.Printf("Run server on %s:%d", api.Host, api.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", api.Host, api.Port), nil)
}

func (api *MediaAPI) ping(rw http.ResponseWriter, r *http.Request) {
	if !api.checkMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	api.setCorsHeaders(&rw)
	fmt.Fprint(rw, "pong")
}

func (api *MediaAPI) signIn(rw http.ResponseWriter, r *http.Request) {
	if !api.checkMethod(r.Method, http.MethodPost) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	api.setCorsHeaders(&rw)
	var creds domain.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	token, err := jwt.GenerateNewToken(creds, api.TokenLiveTime)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	response := domain.Token{
		Token: token,
	}
	json, _ := json.Marshal(response)
	fmt.Fprint(rw, string(json))
	return
}

func (api *MediaAPI) upload(rw http.ResponseWriter, r *http.Request) {
	if !api.checkMethod(r.Method, http.MethodPost) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	api.setCorsHeaders(&rw)
	token := r.Header.Get("Token")
	if token == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	code := api.authorization(token)
	if code != http.StatusOK {
		rw.WriteHeader(code)
		return
	}
	if r.Method != http.MethodPost {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fileBytes, handler, err := r.FormFile("file")
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}
	fileExtension := api.getFileExtension(handler.Filename)
	fileName := api.encodeFileName(handler.Filename, fileExtension)
	fullDataStorageDestination := api.getFullFilePath(fileExtension)

	if _, err := os.Stat(fmt.Sprintf("%s/%s", fullDataStorageDestination, fileName)); err == nil {
		response := domain.FileAlreadyExistError{
			What:     "File already exist!",
			FileName: fileName,
		}
		json, err := json.Marshal(response)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, string(json))
		return
	}
	defer fileBytes.Close()
	os.MkdirAll(fullDataStorageDestination, os.ModePerm)
	data, err := ioutil.ReadAll(fileBytes)
	file, err := os.Create(fmt.Sprintf("%s/%s", fullDataStorageDestination, fileName))
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = file.Write(data)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()
	fmt.Fprint(rw, fileName)
}

func (api *MediaAPI) delete(rw http.ResponseWriter, r *http.Request) {
	if !api.checkMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	api.setCorsHeaders(&rw)
	token := r.Header.Get("Token")
	if token == "" {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}
	code := api.authorization(token)
	if code != http.StatusOK {
		rw.WriteHeader(code)
		return
	}
	if r.Method != http.MethodGet {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	fileName := r.URL.Query().Get("name")
	fileExtension := api.getFileExtension(fileName)
	fullDataStorageDestination := api.getFullFilePath(fileExtension)
	if os.Remove(fmt.Sprintf("%s/%s", fullDataStorageDestination, fileName)) != nil {
		response := domain.Error{
			What: "File not found!",
		}
		json, err := json.Marshal(response)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			return
		}
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, string(json))
		return
	}
	os.Remove(fullDataStorageDestination)
	return
}

func (api *MediaAPI) getFileNamesWithDates(rw http.ResponseWriter, r *http.Request) {
	if !api.checkMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	api.setCorsHeaders(&rw)
	extenstion := r.URL.Query().Get("extension")
	extensionDirectory := api.getFullFilePath(extenstion)
	filesDirectory, err := ioutil.ReadDir(extensionDirectory)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	var filesList []domain.File
	for _, file := range filesDirectory {
		filesList = append(filesList, domain.File{Name: file.Name(), ModTime: file.ModTime()})
	}
	json, err := json.Marshal(filesList)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	fmt.Fprint(rw, string(json))
}

func (api *MediaAPI) getAvailableExtension(rw http.ResponseWriter, r *http.Request) {
	if !api.checkMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	api.setCorsHeaders(&rw)
	path := fmt.Sprintf("%s/%s", api.RootPath, api.DataStorageRoute)
	filesDirectory, err := ioutil.ReadDir(path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	var extensionsList []domain.Extension
	for _, file := range filesDirectory {
		extensionsList = append(extensionsList, domain.Extension{Name: file.Name(), IsDir: file.IsDir()})
	}
	json, err := json.Marshal(extensionsList)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	fmt.Fprint(rw, string(json))
}

func (api *MediaAPI) getDataByUrl(rw http.ResponseWriter, r *http.Request) {
	if !api.checkMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	api.setCorsHeaders(&rw)
	fileName := r.URL.RequestURI()[len(fmt.Sprintf("%s%s", api.Prefix, "/data/")):]
	fileExtension := api.getFileExtension(fileName)
	fullDataStorageDestination := api.getFullFilePath(fileExtension)
	http.ServeFile(rw, r, fmt.Sprintf("%s/%s", fullDataStorageDestination, fileName))
}
