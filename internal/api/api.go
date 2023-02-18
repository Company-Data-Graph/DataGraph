package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"media-server/internal/jwt"
	"media-server/internal/models"
	filesystem "media-server/internal/osProvider"
	"media-server/internal/utils"

	"net/http"
)

type MediaAPI struct {
	Host             string
	Port             int
	Prefix           string
	TokenLiveTime    int
	RootPath         string
	DataStorageRoute string
}

// Create new MediaAPI
func NewMediaAPI(config *models.MediaAPIConfig) *MediaAPI {
	api := &MediaAPI{Host: config.Host, Port: config.Port, Prefix: config.Prefix, TokenLiveTime: config.TokenLiveTime, RootPath: config.StorageRootPath, DataStorageRoute: config.DataStorageRoute}
	jwt.Users = make(map[string]string)
	jwt.Users["admin"] = config.AdminPass
	return api
}

// Run media-server api
func (api *MediaAPI) Run() {
	// Setup api handlers
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/ping/"), api.ping)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/data/"), api.getDataByUrl)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/signin"), api.signIn)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/upload"), api.upload)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/delete"), api.delete)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/dir"), api.getFileNamesWithDates)
	http.HandleFunc(fmt.Sprintf("%s%s", api.Prefix, "/extensions"), api.getAvailableExtension)

	// Start media-server api
	log.Printf("Run server on %s:%d", api.Host, api.Port)
	http.ListenAndServe(fmt.Sprintf("%s:%d", api.Host, api.Port), nil)
}

// Default get ping for check server status
func (api *MediaAPI) ping(rw http.ResponseWriter, r *http.Request) {
	if !utils.CheckMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	utils.SetCorsHeaders(&rw)
	fmt.Fprint(rw, "pong")
}

// Authorization handler
func (api *MediaAPI) signIn(rw http.ResponseWriter, r *http.Request) {
	// Check method allowed
	if !utils.CheckMethod(r.Method, http.MethodPost) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set cors
	utils.SetCorsHeaders(&rw)

	// Decode credentials for authorization
	var creds models.Credentials
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Generate token if credentials is equal
	token, err := jwt.GenerateNewToken(creds, api.TokenLiveTime)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	response := models.Token{
		Token: token,
	}
	json, _ := json.Marshal(response)
	fmt.Fprint(rw, string(json))
	return
}

// Uploading file into media-server (into data storage dir)
func (api *MediaAPI) upload(rw http.ResponseWriter, r *http.Request) {
	// Check method allowed
	if !utils.CheckMethod(r.Method, http.MethodPost) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	// Set cors
	utils.SetCorsHeaders(&rw)

	// Authorization with token
	token := r.Header.Get("Token")
	if jwt.Auth(token) != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Get fileBytes from form and defer buffer closing
	fileBytes, handler, err := r.FormFile("file")
	defer fileBytes.Close()
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		return
	}

	// Calculate file meta data
	fileExtension := utils.GetFileExtension(handler.Filename)
	fileName := utils.EncodeFileName(handler.Filename, fileExtension)
	fullDataStorageDestination := utils.GetFullFilePath(api.RootPath, api.DataStorageRoute, fileExtension)

	// Check file exist
	if filesystem.CheckFileExist(fullDataStorageDestination, fileName) {
		response := models.FileAlreadyExistError{
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

	// Create extension dir if not exist
	filesystem.CreateDir(fullDataStorageDestination)
	// Reading buffer and create file in dir
	data, err := ioutil.ReadAll(fileBytes)
	if filesystem.CreateFile(fileName, fullDataStorageDestination, data) != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Fprint(rw, fileName)
}

func (api *MediaAPI) delete(rw http.ResponseWriter, r *http.Request) {
	// Check method allowed
	if !utils.CheckMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	// Set cors
	utils.SetCorsHeaders(&rw)

	// Authorization with token
	token := r.Header.Get("Token")
	if jwt.Auth(token) != nil {
		rw.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Getting filename from query params and calculate file metadata
	fileName := r.URL.Query().Get("name")
	fileExtension := utils.GetFileExtension(fileName)
	fullDataStorageDestination := utils.GetFullFilePath(api.RootPath, api.DataStorageRoute, fileExtension)

	// Check file exist
	if !filesystem.CheckFileExist(fullDataStorageDestination, fileName) {
		response := models.FileAlreadyExistError{
			What:     "File not exist!",
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

	if err := filesystem.DeleteFile(fileName, fullDataStorageDestination); err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func (api *MediaAPI) getFileNamesWithDates(rw http.ResponseWriter, r *http.Request) {
	if !utils.CheckMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	utils.SetCorsHeaders(&rw)
	extenstion := r.URL.Query().Get("extension")
	extensionDirectory := utils.GetFullFilePath(api.RootPath, api.DataStorageRoute, extenstion)
	filesDirectory, err := ioutil.ReadDir(extensionDirectory)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	var filesList []models.File
	for _, file := range filesDirectory {
		filesList = append(filesList, models.File{Name: file.Name(), ModTime: file.ModTime()})
	}
	json, err := json.Marshal(filesList)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	fmt.Fprint(rw, string(json))
}

func (api *MediaAPI) getAvailableExtension(rw http.ResponseWriter, r *http.Request) {
	if !utils.CheckMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	utils.SetCorsHeaders(&rw)
	path := fmt.Sprintf("%s/%s", api.RootPath, api.DataStorageRoute)
	filesDirectory, err := ioutil.ReadDir(path)
	if err != nil {
		rw.WriteHeader(http.StatusNotFound)
		return
	}
	var extensionsList []models.Extension
	for _, file := range filesDirectory {
		extensionsList = append(extensionsList, models.Extension{Name: file.Name(), IsDir: file.IsDir()})
	}
	json, err := json.Marshal(extensionsList)
	if err != nil {
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	fmt.Fprint(rw, string(json))
}

func (api *MediaAPI) getDataByUrl(rw http.ResponseWriter, r *http.Request) {
	if !utils.CheckMethod(r.Method, http.MethodGet) {
		rw.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	utils.SetCorsHeaders(&rw)
	fileName := r.URL.RequestURI()[len(fmt.Sprintf("%s%s", api.Prefix, "/data/")):]
	fileExtension := utils.GetFileExtension(fileName)
	fullDataStorageDestination := utils.GetFullFilePath(api.RootPath, api.DataStorageRoute, fileExtension)
	http.ServeFile(rw, r, fmt.Sprintf("%s/%s", fullDataStorageDestination, fileName))
}
