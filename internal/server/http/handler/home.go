package handler

import (
	"net/http"
)

func Home() func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		http.ServeFile(response, request, "/Users/anthony/go/src/github.com/instabledesign/mailcatcher/asset/index.html")
		return
	}
}
