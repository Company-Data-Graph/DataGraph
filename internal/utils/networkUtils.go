package utils

import "net/http"

func SetCorsHeaders(rw *http.ResponseWriter) {
	(*rw).Header().Set("Access-Control-Allow-Origin", "*")
	(*rw).Header().Set("Access-Control-Allow-Headers", "Content-Type")
}

func CheckMethod(getted string, required string) bool {
	if required != getted {
		return false
	}
	return true
}
