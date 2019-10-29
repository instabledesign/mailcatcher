package handler

import (
	"encoding/json"
	"net/http"

	"github.com/gol4ng/mailcatcher/internal"
)

func Mail(mailStorage *internal.MailStorage) func(http.ResponseWriter, *http.Request) {
	return func(response http.ResponseWriter, request *http.Request) {
		data := []byte("[]")
		if mails := mailStorage.GetMails(); mails != nil {
			j, err := json.Marshal(mailStorage.GetMails())
			if err != nil {
				response.WriteHeader(http.StatusInternalServerError)
				return
			}
			data = j
		}
		response.Header().Set("Content-type", "application/json")
		response.Write(data)
	}
}
