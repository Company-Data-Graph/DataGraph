package tools

import (
	"fmt"
	"net/http"
)

func SetHttpErrorAndResponse(rw *http.ResponseWriter, code int, message string) {
	(*rw).WriteHeader(code)
	fmt.Fprintf((*rw), "\"error\":%s", message)
	return
}
