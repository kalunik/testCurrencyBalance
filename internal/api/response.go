package api

import (
	"github.com/kalunik/testCurrencyBalance/pkg/log"
	"net/http"
)

func responseError(status int, message string, err error, w http.ResponseWriter, log log.Logger) {
	log.Errorf("%s: %w", message, err)
	http.Error(w, message, status)
}
